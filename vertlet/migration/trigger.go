package migration

import (
	"log"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
)

func migrationTriggered(gceService *compute.Service) error {
	start := time.Now()

	// Signal that the migration has begun.
	err := instances.SetInstanceState(instances.StateMigrating, hostname, gceService)
	if err != nil {
		return err
	}

	// Tell the remote Vertlet to begin.

	log.Printf("Request(Migration Triggered) took %s", time.Since(start))
	return nil
}
