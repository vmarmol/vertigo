package migration

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
)

func handleMigrationStart(w http.ResponseWriter, r *url.URL, gceService *compute.Service) error {
	start := time.Now()

	// Start the migration.
	err := instances.SetInstanceState(instances.StateWarmingUp, hostname, gceService)
	if err != nil {
		return err
	}

	// Get the image to start.
	// TODO(vmarmol):

	// Start the container.
	// TODO(vmarmol)

	// We're done
	err = instances.SetInstanceState(instances.StateOk, hostname, gceService)
	if err != nil {
		return err
	}

	log.Printf("Request(%s) took %s", MigrationStartHandler, time.Since(start))
	return nil
}
