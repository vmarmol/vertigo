package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var argPort = flag.Int("port", 80, "port to listen")
var delay = flag.Int("delay", 250, "number of milliseconds to be busy.")
var stateFile = flag.String("statefile", "/tmp/state", "where to store the state of the container")

func main() {
	flag.Parse()
	startTime := time.Now()

	// Load state from statefile if it exists.
	if _, err := os.Stat(*stateFile); err == nil {
		out, err := ioutil.ReadFile(*stateFile)
		if err != nil {
			log.Fatalf("failed to read checkpoint at %q: %s", *stateFile, err)
		}
		err = startTime.UnmarshalBinary(out)
		if err != nil {
			log.Fatalf("failed to parseh checkpoint at %q: %s", *stateFile, err)
		}
	}

	// Register signal handler to checkpoint our state when we're being migrated.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		out, err := startTime.MarshalBinary()
		if err != nil {
			log.Fatalf("Failed to checkpoint: %s", err)
		}
		err = ioutil.WriteFile(*stateFile, out, 0644)
		if err != nil {
			log.Fatalf("Failed write checkpoint to %q: %s", *stateFile, err)
		}
		os.Exit(0)
	}()

	// Register our work handler, to burn CPU.
	nanoDelay := int64(*delay * 1000 * 1000) // Milli to Nano
	http.HandleFunc("/burn", func(w http.ResponseWriter, r *http.Request) {
		reqTime := time.Now()
		timeSince := time.Since(reqTime)
		for i := 0; timeSince.Nanoseconds() < nanoDelay; i++ {
			f := float64(i)
			n := float64(i + 2)
			f = f / n
			n = n / f
			timeSince = time.Since(reqTime)
		}
		fmt.Fprintf(w, time.Since(startTime).String())
	})
	// Register a ping that does not burn cpu. Used by the health checker.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, time.Since(startTime).String())
	})

	addr := fmt.Sprintf(":%v", *argPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
