package migration

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"code.google.com/p/google-api-go-client/compute/v1"
)

var MigrationStartHandler = "/migration/start"
var MigrationMigrateHandler = "/migration/migrate"
var MigrationDoneHandler = "/migration/done"

var hostname = ""

func init() {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	hostname = name
}

func RegisterHandlers(gceService *compute.Service) {
	http.HandleFunc(MigrationStartHandler, func(w http.ResponseWriter, r *http.Request) {
		err := handleMigrationStart(w, r.URL, gceService)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})

	http.HandleFunc(MigrationDoneHandler, func(w http.ResponseWriter, r *http.Request) {
		err := handleMigrationDone(w, r.URL, gceService)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})

	http.HandleFunc(MigrationMigrateHandler, func(w http.ResponseWriter, r *http.Request) {
		err := handleMigration("TODO", w, gceService)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})
}
