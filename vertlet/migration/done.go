package migration

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
)

func handleMigrationDone(w http.ResponseWriter, r *url.URL, gceService *compute.Service) error {
	start := time.Now()

	// "Turn-down" the instance, clear the Vertigo state.
	err := instances.ClearVertigoState(hostname, gceService)
	if err != nil {
		return err
	}

	log.Printf("Request(%s) took %s", MigrationDoneHandler, time.Since(start))
	return nil
}
