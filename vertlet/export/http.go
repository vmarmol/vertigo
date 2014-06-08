package export

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"strings"
)

var ExportContainerHandler = "/export/"

type taskExport struct {
	export ContainerExporter
}

func (self *taskExport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	args := strings.Split(r.URL.Path, "/")
	switch strings.ToUpper(r.Method) {
	case "GET":
		id := args[len(args)-1]
		fmt.Printf("export %v\n", id)
		err := self.export.Export(id, w)
		if err != nil {
			log.Printf("Error when exporting %v: %v\n", id, err)
		}
	}
}

func Register() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	if currentUser.Username != "root" {
		return fmt.Errorf("must be root!")
	}
	dockerExport, err := NewDockerExporter()
	if err != nil {
		return err
	}
	export := &taskExport{
		export: dockerExport,
	}
	http.Handle(ExportContainerHandler, export)
	return nil
}
