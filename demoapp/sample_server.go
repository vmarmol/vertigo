package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var argPort = flag.Int("port", 80, "port to listen")
var delay = flag.Int("delay", 250, "number of milliseconds to be busy.")
var stateFile = flag.String("statefile", "/tmp/state", "where to store the state of the container")

var percentBurnLock sync.Mutex
var percentBurn = 20
var nsInSecond = int64(1000000000)

func burner() {
	burn := 0
	for true {
		func() {
			percentBurnLock.Lock()
			defer percentBurnLock.Unlock()
			burn = percentBurn
		}()

		burnNs := int64(int64(burn) * (nsInSecond) / int64(100))
		sleepNs := nsInSecond - burnNs

		reqTime := time.Now()
		timeSince := time.Since(reqTime)
		for i := 0; timeSince.Nanoseconds() < burnNs; i++ {
			f := float64(i)
			n := float64(i + 2)
			f = f / n
			n = n / f
			timeSince = time.Since(reqTime)
		}

		time.Sleep(time.Duration(sleepNs))
	}
}

type checkpoint struct {
	StartTime time.Time `json:"start_time"`
	Percent   int       `json:"percent"`
}

func main() {
	flag.Parse()
	startTime := time.Now()

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// Load state from statefile if it exists.
	if _, err := os.Stat(*stateFile); err == nil {
		out, err := ioutil.ReadFile(*stateFile)
		if err != nil {
			log.Fatalf("failed to read checkpoint at %q: %s", *stateFile, err)
		}
		var ck checkpoint
		err = json.Unmarshal(out, &ck)
		if err != nil {
			log.Fatalf("failed to parseh checkpoint at %q: %s", *stateFile, err)
		}
		startTime = ck.StartTime
		percentBurn = ck.Percent
	}

	// Register signal handler to checkpoint our state when we're being migrated.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ck := checkpoint{
			StartTime: startTime,
			Percent:   percentBurn,
		}
		out, err := json.Marshal(ck)
		if err != nil {
			log.Fatalf("Failed to checkpoint: %s", err)
		}
		err = ioutil.WriteFile(*stateFile, out, 0644)
		if err != nil {
			log.Fatalf("Failed write checkpoint to %q: %s", *stateFile, err)
		}
		os.Exit(0)
	}()

	log.Printf("Start time %s and percent burn %d", time.Since(startTime), percentBurn)

	// Start a burner thread per CPU.
	for i := 0; i < runtime.NumCPU(); i++ {
		go burner()
	}

	// Register our work handler, to burn CPU.
	resource := "/burn/"
	http.HandleFunc(resource, func(w http.ResponseWriter, r *http.Request) {
		units, err := strconv.ParseInt(r.URL.Path[len(resource):], 0, 32)
		if err != nil {
			units = 0
		}
		// 4 units per core
		requestPercent := int(units * int64(100) / int64(runtime.NumCPU()*4))
		log.Printf("Requested: %d%%", requestPercent)

		func() {
			percentBurnLock.Lock()
			defer percentBurnLock.Unlock()
			percentBurn = int(requestPercent)
		}()

		fmt.Fprintf(w, time.Since(startTime).String())
	})
	// Register a ping that does not burn cpu. Used by the health checker.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, time.Since(startTime).String())
	})

	addr := fmt.Sprintf(":%v", *argPort)
	log.Printf("Serving on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
