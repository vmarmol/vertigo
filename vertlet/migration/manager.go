package migration

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/fsouza/go-dockerclient"
	"github.com/vmarmol/vertigo/vertlet/monitor"
)

type MigrationHandler struct {
	port             int
	gceService       *compute.Service
	hostname         string
	dockerClient     *docker.Client
	containerTracker monitor.ContainerTracker
}

var argCadvisorUrl = flag.String("cadvisor", "http://localhost:5000", "cadvisor address")
var argCpuLowThreshold = flag.Float64("low", 0.1, "low threshold for cpu")
var argCpuHighThreshold = flag.Float64("high", 0.9, "high threshold for cpu")

func NewMigrationHandler(port int, gceService *compute.Service) (*MigrationHandler, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return nil, err
	}

	sigChan := make(chan *monitor.MonitorSignal)

	tracker, err := monitor.StartDockerMonitor(
		*argCadvisorUrl,
		*argCpuLowThreshold,
		*argCpuHighThreshold,
		sigChan,
	)
	if err != nil {
		log.Fatal(err)
	}

	mig := &MigrationHandler{
		port:             port,
		gceService:       gceService,
		hostname:         name,
		dockerClient:     client,
		containerTracker: tracker,
	}

	go func() {
		for sig := range sigChan {
			id := path.Base(sig.ContainerName)
			switch sig.MoveDst {
			case monitor.DST_LOWER:
				mig.Migrate(id, false)
			case monitor.DST_HIGHER:
				mig.Migrate(id, true)
			}
		}
	}()
	return mig, nil

}

var MigrationStartHandler = "/migration/start"
var MigrationMigrateHandler = "/migration/migrate/"
var TrackContainer = "/track/"

func (self *MigrationHandler) RegisterHandlers() {
	http.HandleFunc(MigrationStartHandler, func(w http.ResponseWriter, r *http.Request) {
		err := self.handleMigrationStart(w, r)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})

	http.HandleFunc(TrackContainer, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		id := path.Base(r.URL.Path)
		err := self.containerTracker.TrackContainer(id)
		if err != nil {
			log.Printf("error when tracking: %v")
		}
	})

	http.HandleFunc(MigrationMigrateHandler, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
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
