package migration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
	"github.com/vmarmol/vertigo/pulet"
)

type MigrationRequest struct {
	Container string
	Host      string
	Port      int
	Command   []string
}

func handleMigrationStart(w http.ResponseWriter, r *http.Request, gceService *compute.Service) error {
	start := time.Now()

	// Get request from the body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var request MigrationRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		return err
	}

	// Start the migration.
	err = instances.SetInstanceState(instances.StateWarmingUp, hostname, gceService)
	if err != nil {
		return err
	}

	pul := pulet.NewPulet()
	importSpec := &pulet.ImportSpec{
		SourceHost: request.Host,
		SourcePort: request.Port,
		SourceId:   request.Container,
	}

	// Get the image to start.
	img, err := pul.Import(importSpec)
	if err != nil {
		log.Fatal(err)
	}

	// Start the container.
	err = pul.RunImage(img, nil, request.Command)
	if err != nil {
		log.Fatal(err)
	}

	// We're done.
	err = instances.SetInstanceState(instances.StateOk, hostname, gceService)
	if err != nil {
		return err
	}

	log.Printf("Request(%s) for container %q took %s", MigrationStartHandler, request.Container, time.Since(start))
	return nil
}
