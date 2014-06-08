package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/vmarmol/vertigo/gce"
	"github.com/vmarmol/vertigo/ui/api"
	"github.com/vmarmol/vertigo/ui/static"
)

var argPort = flag.Int("port", 8080, "port to listen")

func main() {
	gceService, err := gce.NewCompute()
	if err != nil {
		log.Fatal(err)
	}

	// Handler for static content.
	http.HandleFunc(static.StaticResource, func(w http.ResponseWriter, r *http.Request) {
		err := static.HandleRequest(w, r.URL)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})

	// Handler for instance information.
	http.HandleFunc(api.InstancesResource, func(w http.ResponseWriter, r *http.Request) {
		err := api.GetInstances(gceService, w)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	})

	// Redirect / to /static/index.html.
	http.Handle("/", http.RedirectHandler(path.Join("/", static.StaticResource, "/index.html"), http.StatusTemporaryRedirect))

	addr := fmt.Sprintf(":%v", *argPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
