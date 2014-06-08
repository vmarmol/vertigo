package migration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/vmarmol/vertigo/instances"
	"github.com/vmarmol/vertigo/pulet"
)

type MigrationRequest struct {
	Container string
	Host      string
	Port      int
	Command   []string
}

func (self *MigrationHandler) handleMigrationStart(w http.ResponseWriter, r *http.Request) error {
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
	log.Printf("Warming up...")
	err = instances.SetInstanceState(instances.StateWarmingUp, self.hostname, self.gceService)
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
	log.Printf("Importing container...")
	img, err := pul.Import(importSpec)
	if err != nil {
		log.Fatal(err)
	}

	// Start the container.
	log.Printf("Running image...")
	// FIXME(monnand): Wrong way!
	args := []string{"-p", "0.0.0.0:80:80"}
	err = pul.RunImage(img, args, request.Command)
	if err != nil {
		log.Fatalf("Error running image: %s", err)
	}

	// We're done.
	log.Printf("Migration complete!")
	err = instances.SetInstanceState(instances.StateOk, self.hostname, self.gceService)
	if err != nil {
		return err
	}

	err = self.containerTracker.TrackContainer(request.Container)
	if err != nil {
		return err
	}

	log.Printf("Request(%s) for container %q took %s", MigrationStartHandler, request.Container, time.Since(start))
	return nil
}
