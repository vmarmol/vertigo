package vertigo

import (
	"fmt"
	"net/http"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/vmarmol/vertigo/instances"
)

var VertigoAddHandler = "/vertigo/add/"
var VertigoRemoveHandler = "/vertigo/remove/"

func getInstance(handler, url string) string {
	if len(handler) >= len(url) {
		return ""
	}

	return url[len(handler):]
}

func RegisterHandlers(gceService *compute.Service) {
	http.HandleFunc(VertigoAddHandler, func(w http.ResponseWriter, r *http.Request) {
		instance := getInstance(VertigoAddHandler, r.URL.Path)
		if instance == "" {
			fmt.Fprintf(w, "no instance specified")
		}
		err := instances.SetInstanceState(instances.StateOk, instance, gceService)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})
	http.HandleFunc(VertigoRemoveHandler, func(w http.ResponseWriter, r *http.Request) {
		instance := getInstance(VertigoRemoveHandler, r.URL.Path)
		if instance == "" {
			fmt.Fprintf(w, "no instance specified")
		}
		err := instances.ClearVertigoState(instance, gceService)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})
}
