package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/vmarmol/vertigo/gce"
)

var argPort = flag.Int("port", 8080, "port to listen")

func main() {
	compute, err := gce.NewCompute()
	if err != nil {
		log.Fatal(err)
	}
	ins, err := compute.Instances.List("lmctfy-prod", "us-central1-a").Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Instances: %v", ins)

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	addr := fmt.Sprintf(":%v", *argPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
