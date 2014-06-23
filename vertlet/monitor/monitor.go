package monitor

import (
	"log"
	"time"

	"github.com/google/cadvisor/client"
	"github.com/google/cadvisor/info"
)

type ContainerMonitor struct {
	containerName string
	client        *cadvisor.Client
	stop          chan bool
	numCores      int
	createdAt     time.Time
}

const (
	DST_HIGHER = iota
	DST_LOWER
)

type MonitorSignal struct {
	MoveDst       int
	ContainerName string
}

func NewContainerMonitor(
	cadvisorUrl string,
	containerName string,
	cpuLowThreshold,
	cpuHighThreshold float64,
	sigChan chan<- *MonitorSignal,
) (*ContainerMonitor, error) {
	c, err := cadvisor.NewClient(cadvisorUrl)
	if err != nil {
		return nil, err
	}

	minfo, err := c.MachineInfo()
	m := &ContainerMonitor{
		containerName: containerName,
		client:        c,
		stop:          make(chan bool),
		numCores:      minfo.NumCores,
		createdAt:     time.Now(),
	}

	go m.checkContainer(3*time.Second, func(util float64) {
		sig := &MonitorSignal{
			ContainerName: containerName,
		}
		log.Printf("container %v, cpu util: %v; [%v,%v]\n", containerName, util, cpuLowThreshold, cpuHighThreshold)
		if util > cpuHighThreshold {
			sig.MoveDst = DST_HIGHER
			sigChan <- sig
		} else if util < cpuLowThreshold {
			sig.MoveDst = DST_LOWER
			sigChan <- sig
		}
	})
	return m, nil
}

func getLatestStats(cinfo *info.ContainerInfo) *info.ContainerStats {
	var latest time.Time
	var ret *info.ContainerStats

	for _, s := range cinfo.Stats {
		if s.Timestamp.After(latest) {
			latest = s.Timestamp
			ret = s
		}
	}
	return ret
}

func (self *ContainerMonitor) checkContainer(
	sleepDuration time.Duration,
	callback func(cpuUtil float64),
) {
	var prevStats *info.ContainerStats
	// Let the container warm up
	time.Sleep(5 * time.Second)
	for {
		cinfo, err := self.client.ContainerInfo(self.containerName)
		if err != nil {
			return
		}
		stats := getLatestStats(cinfo)
		if prevStats != nil {
			CpuDiff := stats.Cpu.Usage.Total - prevStats.Cpu.Usage.Total
			util := float64(CpuDiff) / float64(sleepDuration.Nanoseconds()*int64(self.numCores))
			callback(util)
		}
		prevStats = stats
		select {
		case <-self.stop:
			return
		case <-time.After(sleepDuration):
		}
	}
}

func (self *ContainerMonitor) Stop() {
	close(self.stop)
}
