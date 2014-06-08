package migration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
)

func handleMigration(request MigrationRequest, remoteVertlet string, gceService *compute.Service) error {
	start := time.Now()

	// Signal that the migration has begun.
	err := instances.SetInstanceState(instances.StateMigrating, hostname, gceService)
	if err != nil {
		return err
	}

	// Tell the remote Vertlet to migrate.
	encodedRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post("http://"+remoteVertlet+MigrationStartHandler, "application/json", bytes.NewReader(encodedRequest))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// The remote Vertlet finished, "turn-down" the instance, clear the Vertigo state.
	err = instances.ClearVertigoState(hostname, gceService)
	if err != nil {
		return err
	}

	// TODO(vmarmol): Do we rm the container?

	log.Printf("Request(%s) took %s", MigrationMigrateHandler, time.Since(start))
	return nil
}
