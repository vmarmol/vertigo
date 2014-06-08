package gce

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// XXX(monnand): We know we should use oauth. But this is a hackathon.
type gcutilVmManager struct {
	gcutilPath string
}

const (
	gcutil_OP_ADD = iota
	gcutil_OP_DEL
	gcutil_OP_ADD_TO_POOL
)

type gcutilOperation int

func (self gcutilOperation) String() string {
	switch self {
	case gcutil_OP_ADD:
		return "addinstance"
	case gcutil_OP_DEL:
		return "deleteinstance"
	case gcutil_OP_ADD_TO_POOL:
		return "addtargetpoolinstance"
	}
	return ""
}

type gcutlCmdParams struct {
	Name        string
	Op          gcutilOperation
	GcutilPath  string
	Zone        string
	MachineType string
	Image       string
	TargetPool  string
}

func (self *gcutlCmdParams) fillDefault() {
	if self.Name == "" {
		self.Name = fmt.Sprintf("vertigo-%v", randomUniqString())
	}
	if self.Zone == "" {
		self.Zone = "us-central1-a"
	}
	if self.MachineType == "" {
		self.MachineType = "n1-standard-2"
	}
	if self.Image == "" {
		self.Image = "ubuntu-trusty"
	}
	if self.TargetPool == "" {
		self.TargetPool = "vertigo-lb-pool"
	}
}

func (self *gcutlCmdParams) ToParamList() []string {
	self.fillDefault()
	ret := make([]string, 0, 6)
	ret = append(ret, self.Op.String())
	switch self.Op {
	case gcutil_OP_ADD:
		ret = append(ret, fmt.Sprintf("--zone=%v", self.Zone))
		ret = append(ret, fmt.Sprintf("--machine_type=%v", self.MachineType))
		ret = append(ret, fmt.Sprintf("--image=%v", self.Image))
		ret = append(ret, self.Name)
	case gcutil_OP_DEL:
		ret = append(ret, "-f")
		ret = append(ret, "--nodelete_boot_pd")
		ret = append(ret, self.Name)
	case gcutil_OP_ADD_TO_POOL:
		ret = append(ret, fmt.Sprintf("--instances=%v", self.Name))
		ret = append(ret, "--region=us-central1")
		ret = append(ret, self.TargetPool)
	}
	return ret
}

func specToAddInstanceCmd(spec *VirtualMachineSpec) *gcutlCmdParams {
	ret := &gcutlCmdParams{
		Op:    gcutil_OP_ADD,
		Name:  spec.Name,
		Image: spec.Image,
	}
	return ret
}

func infoToDelInstanceCmd(info *VirtualMachineInfo) *gcutlCmdParams {
	return &gcutlCmdParams{
		Op:   gcutil_OP_DEL,
		Name: info.Name,
	}
}

func specToAddToTargetPool(spec *VirtualMachineSpec) *gcutlCmdParams {
	return &gcutlCmdParams{
		Op:    gcutil_OP_ADD_TO_POOL,
		Name:  spec.Name,
		Image: spec.Image,
	}
}

func (self *gcutilVmManager) runGcutil(params *gcutlCmdParams, timeout time.Duration) error {
	ch := make(chan bool)
	var err error
	go func() {
		paramList := params.ToParamList()
		gcutil := self.gcutilPath
		if self.gcutilPath == "" {
			gcutil = "gcutil"
		}
		fmt.Printf("%v %v\n", gcutil, strings.Join(paramList, " "))
		cmd := exec.Command(gcutil, paramList...)
		err = cmd.Run()
		ch <- true
	}()
	if timeout.Seconds() > 1.0 {
		select {
		case <-time.After(timeout):
		case <-ch:
		}
	} else {
		<-ch
	}
	if err != nil {
		return fmt.Errorf("unable to use gcutil: %v", err)
	}
	return nil
}

func (self *gcutilVmManager) NewMachine(spec *VirtualMachineSpec) (*VirtualMachineInfo, error) {
	params := specToAddInstanceCmd(spec)
	err := self.runGcutil(params, 0*time.Second)
	if err != nil {
		return nil, err
	}
	info := &VirtualMachineInfo{
		Name: params.Name,
	}

	params = specToAddToTargetPool(spec)
	params.Name = info.Name
	err = self.runGcutil(params, 30*time.Second)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (self *gcutilVmManager) DelMachine(info *VirtualMachineInfo) error {
	params := infoToDelInstanceCmd(info)
	err := self.runGcutil(params, 0*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func NewGcutilManager(gcutilPath string) (VirtualMachineManager, error) {
	ret := &gcutilVmManager{
		gcutilPath: "gcutil",
	}
	return ret, nil
}
