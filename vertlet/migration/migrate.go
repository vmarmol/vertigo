package migration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/vmarmol/vertigo/instances"
)

func (self *MigrationHandler) Migrate(container string, migrateUp bool) error {
	// Get the command from Docker.
	cont, err := self.dockerClient.InspectContainer(container)
	if err != nil {
		return err
	}
	command := make([]string, 0, len(cont.Args)+1)
	command = append(command, cont.Path)
	command = append(command, cont.Args...)

	request := MigrationRequest{
		Container: container,
		Host:      self.hostname,
		Port:      self.port,
		Command:   command,
	}

	// Find where to migrate to.
	var destination string
	if migrateUp {
		destination, err = instances.GetLargerInstance(self.hostname)
		if err != nil {
			return err
		}
	} else {
		destination, err = instances.GetSmallerInstance(self.hostname)
		if err != nil {
			return err
		}
	}
	destination = fmt.Sprintf("%s:%d", destination, request.Port)

	log.Printf("Migrating container %q to %q", container, destination)

	return self.handleMigration(request, destination)
}

func (self *MigrationHandler) handleMigration(request MigrationRequest, remoteVertlet string) error {
	start := time.Now()

	// Signal that the migration has begun.
	err := instances.SetInstanceState(instances.StateMigrating, self.hostname, self.gceService)
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
	err = instances.ClearVertigoState(self.hostname, self.gceService)
	if err != nil {
		return err
	}

	// TODO(vmarmol): Do we rm the container?

	log.Printf("Request(%s) took %s", MigrationMigrateHandler, time.Since(start))
	return nil
}
