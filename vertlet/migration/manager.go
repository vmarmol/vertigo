package migration

import (
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/google-api-go-client/compute/v1"
)

type MigrationHandler struct {
	port       int
	gceService *compute.Service
	hostname   string
}

func NewMigrationHandler(port int, gceService *compute.Service) (*MigrationHandler, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &MigrationHandler{
		port:       port,
		gceService: gceService,
		hostname:   name,
	}, nil
}

var MigrationStartHandler = "/migration/start"
var MigrationMigrateHandler = "/migration/migrate"

func (self *MigrationHandler) RegisterHandlers() {
	http.HandleFunc(MigrationStartHandler, func(w http.ResponseWriter, r *http.Request) {
		err := self.handleMigrationStart(w, r)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})

	http.HandleFunc(MigrationMigrateHandler, func(w http.ResponseWriter, r *http.Request) {
		err := self.Migrate("c7d4f0543e92", []string{"/bin/sleep", "2m"}, true)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})
}
