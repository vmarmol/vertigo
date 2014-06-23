package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/google/cadvisor/client"
	"github.com/vmarmol/vertigo/instances"
)

var InstancesResource = "/api/instances"

type Instance struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	CpuUsage    int    `json:"cpu_usage"`
	MemoryUsage int    `json:"memory_usage"`
}

type trackedContainer struct {
	Tracked string `json:"tracked"`
}

func getUsage(instance string) (int, int, error) {
	// Ask the Vertlet what the tracked container is.
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/tracked", instance))
	if err != nil {
		return -1, -1, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, err
	}
	var tracked trackedContainer
	err = json.Unmarshal(body, &tracked)
	if err != nil {
		return -1, -1, err
	}
	trackedId := tracked.Tracked

	// Not tracking a container, no usage then.
	if trackedId == "" {
		return 0, 0, nil
	}

	// Get the usage from cAdvisor.
	log.Printf("Instance %q is tracking: %q", instance, trackedId)
	c, err := cadvisor.NewClient(fmt.Sprintf("http://%s:5000/", instance))
	if err != nil {
		return -1, -1, err
	}
	cinfo, err := c.ContainerInfo(trackedId)
	if err != nil {
		return -1, -1, err
	}
	statsLen := len(cinfo.Stats)
	if statsLen < 2 {
		return 0, 0, nil
	}

	// Get the machine info from cAdvisor.
	m, err := c.MachineInfo()
	if err != nil {
		return -1, -1, err
	}
	cpuUsage := cinfo.Stats[statsLen-1].Cpu.Usage.Total - cinfo.Stats[statsLen-2].Cpu.Usage.Total
	return int(cpuUsage*uint64(100)/1000000000) / m.NumCores, int(int64(cinfo.Stats[statsLen-1].Memory.Usage) * int64(100) / m.MemoryCapacity), nil
}

func GetInstances(serv *compute.Service, w http.ResponseWriter) error {
	start := time.Now()

	// Get instances.
	instances, err := instances.GetVertigoInstances(serv)
	if err != nil {
		return err
	}
	output := make([]Instance, 0, len(instances))
	for _, instance := range instances {
		log.Printf("Instance: %q", instance)
		cpu, mem, err := getUsage(instance.Name)
		if err != nil {
			return err
		}

		output = append(output, Instance{
			Name:        instance.Name,
			State:       instance.State,
			CpuUsage:    cpu,
			MemoryUsage: mem,
		})
	}

	// Marshall and return the output.
	out, err := json.Marshal(output)
	if err != nil {
		return err
	}
	w.Write(out)

	log.Printf("Request(%s) took %s", InstancesResource, time.Since(start))
	return nil
}
