package monitor

import (
	"fmt"
	"log"
	"sync"
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
	lock             sync.Mutex
}

type ContainerTracker interface {
	TrackContainer(id string) error
}

func StartDockerMonitor(
	cadvisorUrl string,
	cpuLowThreshold,
	cpuHighThreshold float64,
	sigChan chan<- *MonitorSignal,
) (ContainerTracker, error) {
	c, err := cadvisor.NewClient(cadvisorUrl)
	if err != nil {
		return nil, err
	}
	m := &DockerMonitor{
		client:           c,
		subcontainers:    make(map[string]*ContainerMonitor, 10),
		cadvisorUrl:      cadvisorUrl,
		cpuLowThreshold:  cpuLowThreshold,
		cpuHighThreshold: cpuHighThreshold,
		sigChan:          sigChan,
	}
	// go m.checkDockerContainers()

	return m, err
}

func (self *DockerMonitor) TrackContainer(id string) error {
	cinfo, err := self.client.ContainerInfo("/docker")
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	cpath := fmt.Sprintf("/docker/%v", id)
	self.lock.Lock()
	defer self.lock.Unlock()
	for _, sub := range cinfo.Subcontainers {
		if sub == cpath {
			if _, ok := self.subcontainers[cpath]; ok {
				return nil
			}
			m, err := NewContainerMonitor(
				self.cadvisorUrl,
				cpath,
				self.cpuLowThreshold,
				self.cpuHighThreshold,
				self.sigChan)
			if err != nil {
				return err
			}
			self.subcontainers[cpath] = m
			return nil
		}
	}
	return fmt.Errorf("cannot find container %v", id)
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
