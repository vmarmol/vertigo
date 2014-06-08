package gce

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

type VirtualMachineSpec struct {
	Name        string
	CpuLevel    int
	MemoryLevel int
	Image       string
}

func randomUniqString() string {
	var d [8]byte
	io.ReadFull(rand.Reader, d[:])
	str := hex.EncodeToString(d[:])
	return fmt.Sprintf("%x-%v", time.Now().Unix(), str)
}

func (self *VirtualMachineSpec) GetName() string {
	if self.Name == "" {
		return fmt.Sprintf("vertigo-%v", randomUniqString())
	}
	return self.Name
}

type VirtualMachineInfo struct {
	Name    string
	Address string
}

type VirtualMachineManager interface {
	NewMachine(spec *VirtualMachineSpec) (*VirtualMachineInfo, error)
}
