package export

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"
)

type DockerExporter struct {
	client *docker.Client
	lock   sync.Mutex
}

var endpoint = "unix:///var/run/docker.sock"

type ContainerExporter interface {
	//RunTask(runspec *api.RunSpec) (containerSpec *api.ContainerSpec, err error)
	Export(container string, w io.Writer) error
}

func NewDockerExporter() (ContainerExporter, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &DockerExporter{
		client: client,
	}, nil
}

/*
func (self *DockerTaskManager) pull(image string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	cmd := exec.Command("docker", "pull", image)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
*/

func randomUniqString() string {
	var d [8]byte
	io.ReadFull(rand.Reader, d[:])
	str := hex.EncodeToString(d[:])
	return fmt.Sprintf("%x-%v", time.Now().Unix(), str)
}

/*
func (self *DockerTaskManager) RunTask(runspec *api.RunSpec) (containerSpec *api.ContainerSpec, err error) {
	err = self.pull(runspec.Image)
	if err != nil {
		return
	}
	log.Println("pulled image")

	exposedPorts := make(map[docker.Port]struct{}, len(runspec.Ports))
	portBindings := make(map[docker.Port][]docker.PortBinding, len(runspec.Ports))
	for _, port := range runspec.Ports {
		if port.HostPort != port.ContainerPort {
			err = fmt.Errorf("host port != container port: %+v", port)
			return
		}
		dport := docker.Port(fmt.Sprintf("%v/tcp", port.ContainerPort))
		exposedPorts[dport] = struct{}{}
		portBindings[dport] = []docker.PortBinding{
			docker.PortBinding{
				HostPort: fmt.Sprintf("%v", port.HostPort),
			},
		}
	}

	name := randomUniqString()
	env := make([]string, 0, len(runspec.Env))
	for _, e := range runspec.Env {
		env = append(env, e.String())
	}

	cmd := make([]string, 0, 1+len(runspec.Args))
	cmd = append(cmd, runspec.Cmd)
	if len(runspec.Args) > 0 {
		cmd = append(cmd, runspec.Args...)
	}

	opts := docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image:        runspec.Image,
			ExposedPorts: exposedPorts,
			Env:          env,
			Cmd:          cmd,
		},
	}
	log.Printf("creating container %+v\n", opts)
	container, err := self.client.CreateContainer(opts)
	if err != nil {
		return
	}
	log.Printf("created container %+v\n", container)

	err = self.client.StartContainer(container.ID, &docker.HostConfig{
		PortBindings: portBindings,
	})
	if err != nil {
		return
	}
	log.Printf("started container\n")
	containerSpec = &api.ContainerSpec{
		Id: container.ID,
	}
	return
}
*/

func (self *DockerExporter) StopTask(id string) error {
	return self.client.StopContainer(id, 10)
}

func (self *DockerExporter) Export(id string, w io.Writer) error {
	err := self.StopTask(id)
	if err != nil {
		return err
	}
	opts := docker.ExportContainerOptions{
		ID:           id,
		OutputStream: w,
	}
	err = self.client.ExportContainer(opts)
	return err
}
