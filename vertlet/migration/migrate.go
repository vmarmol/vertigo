package migration

import (
	"log"
	"net/http"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
)

func handleMigration(container string, w http.ResponseWriter, gceService *compute.Service) error {
	start := time.Now()

	// Signal that the migration has begun.
	err := instances.SetInstanceState(instances.StateMigrating, hostname, gceService)
	if err != nil {
		return err
	}

	// Stop the container.
	// TODO(vmarmol):

	// Tell the remote Vertlet to begin.
	// TODO(vmarmol)

	log.Printf("Request(%s) took %s", MigrationMigrateHandler, time.Since(start))
	return nil
}
