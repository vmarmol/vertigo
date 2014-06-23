package monitor

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/cadvisor/client"
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
	StopTrackingContainer(id string) error
	GetTrackedContainer() string
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
	log.Printf("Tracking %q", id)
	cinfo, err := self.client.ContainerInfo("/docker")
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	cpath := fmt.Sprintf("/docker/%v", id)
	self.lock.Lock()
	defer self.lock.Unlock()
	for _, sub := range cinfo.Subcontainers {
		if strings.HasPrefix(sub.Name, cpath) {
			if _, ok := self.subcontainers[sub.Name]; ok {
				return nil
			}
			m, err := NewContainerMonitor(
				self.cadvisorUrl,
				sub.Name,
				self.cpuLowThreshold,
				self.cpuHighThreshold,
				self.sigChan)
			if err != nil {
				return err
			}
			self.subcontainers[sub.Name] = m
			return nil
		}
	}
	return fmt.Errorf("cannot find container %v aka: %q", id, cpath)
}

func (self *DockerMonitor) StopTrackingContainer(id string) error {
	cpath := fmt.Sprintf("/docker/%v", id)
	self.lock.Lock()
	defer self.lock.Unlock()
	if m, ok := self.subcontainers[cpath]; ok {
		m.Stop()
		delete(self.subcontainers, cpath)
	} else {
		return fmt.Errorf("unknown container %v", id)
	}
	return nil
}

func (self *DockerMonitor) GetTrackedContainer() string {
	self.lock.Lock()
	defer self.lock.Unlock()
	for id, _ := range self.subcontainers {
		return id
	}
	return ""
}
