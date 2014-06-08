package migration

import (
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/fsouza/go-dockerclient"
)

type MigrationHandler struct {
	port         int
	gceService   *compute.Service
	hostname     string
	dockerClient *docker.Client
}

func NewMigrationHandler(port int, gceService *compute.Service) (*MigrationHandler, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return nil, err
	}
	return &MigrationHandler{
		port:         port,
		gceService:   gceService,
		hostname:     name,
		dockerClient: client,
	}, nil
}

var MigrationStartHandler = "/migration/start"
var MigrationMigrateHandler = "/migration/migrate/"

func (self *MigrationHandler) RegisterHandlers() {
	http.HandleFunc(MigrationStartHandler, func(w http.ResponseWriter, r *http.Request) {
		err := self.handleMigrationStart(w, r)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})

	http.HandleFunc(MigrationMigrateHandler, func(w http.ResponseWriter, r *http.Request) {
		if len(MigrationMigrateHandler) >= len(r.URL.Path) {
			fmt.Fprintf(w, "Missing container name")
			return
		}
		err := self.Migrate(r.URL.Path[len(MigrationMigrateHandler):], true)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})
}
