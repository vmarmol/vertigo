package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
)

var InstancesResource = "/api/instances"

type Instance struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	CpuUsage    int    `json:"cpu_usage"`
	MemoryUsage int    `json:"memory_usage"`
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
		output = append(output, Instance{
			Name:  instance.Name,
			State: instance.State,
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
