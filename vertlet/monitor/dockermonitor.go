package monitor

import (
	"log"
	"time"

	"github.com/google/lmctfy/cadvisor/client"
)

type DockerMonitor struct {
	client           *cadvisor.Client
	subcontainers    map[string]*ContainerMonitor
	cadvisorUrl      string
	cpuLowThreshold  float64
	cpuHighThreshold float64
	sigChan          chan<- *MonitorSignal
}

func StartDockerMonitor(
	cadvisorUrl string,
	cpuLowThreshold,
	cpuHighThreshold float64,
	sigChan chan<- *MonitorSignal,
) error {
	c, err := cadvisor.NewClient(cadvisorUrl)
	if err != nil {
		return err
	}
	m := &DockerMonitor{
		client:           c,
		subcontainers:    make(map[string]*ContainerMonitor, 10),
		cadvisorUrl:      cadvisorUrl,
		cpuLowThreshold:  cpuLowThreshold,
		cpuHighThreshold: cpuHighThreshold,
		sigChan:          sigChan,
	}
	go m.checkDockerContainers()

	return nil
}

func (self *DockerMonitor) addSubContainers(containers []string) {
	for _, c := range containers {
		if _, ok := self.subcontainers[c]; !ok {
			m, err := NewContainerMonitor(self.cadvisorUrl, c, self.cpuLowThreshold, self.cpuHighThreshold, self.sigChan)
			if err != nil {
				log.Printf("unable to create sub container monitor: %v", err)
			}
			self.subcontainers[c] = m
		}
	}
}

func (self *DockerMonitor) checkDockerContainers() {
	for {
		cinfo, err := self.client.ContainerInfo("/docker")
		if err != nil {
			log.Printf("error: %v", err)
		}
		self.addSubContainers(cinfo.Subcontainers)
		time.Sleep(1500 * time.Millisecond)
	}
}
