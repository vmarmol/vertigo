package api

import "fmt"

type Port struct {
	HostPort      int `json:"host_port"`
	ContainerPort int `json:"container_port"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (self Env) String() string {
	return fmt.Sprintf("%v=%v", self.Name, self.Value)
}

type RunSpec struct {
	Image string   `json:"image"`
	Cmd   string   `json:"cmd"`
	Args  []string `json:"args,omitempty"`
	Ports []Port   `json:"ports,omitempty"`
	Env   []Env    `json:"env,omitempty"`
}

type ContainerSpec struct {
	Id string `json:"id"`
}
