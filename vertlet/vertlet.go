package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/vmarmol/vertigo/gce"
	"github.com/vmarmol/vertigo/vertlet/export"
	"github.com/vmarmol/vertigo/vertlet/migration"
	"github.com/vmarmol/vertigo/vertlet/monitor"
)

var argPort = flag.Int("port", 8080, "port to listen")
var argCadvisorUrl = flag.String("cadvisor", "http://localhost:5000", "cadvisor address")
var argCpuLowThreshold = flag.Float64("low", 0.1, "low threshold for cpu")
var argCpuHighThreshold = flag.Float64("high", 0.9, "high threshold for cpu")

func main() {
	flag.Parse()
	gceService, err := gce.NewCompute()
	if err != nil {
		log.Fatal(err)
	}

	err = export.Register()
	if err != nil {
		log.Fatal(err)
	}
	mig, err := migration.NewMigrationHandler(*argPort, gceService)
	if err != nil {
		log.Fatal(err)
	}
	mig.RegisterHandlers()

	sigChan := make(chan *monitor.MonitorSignal)

	err = monitor.StartDockerMonitor(
		*argCadvisorUrl,
		*argCpuLowThreshold,
		*argCpuHighThreshold,
		sigChan,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for sig := range sigChan {
			switch sig.MoveDst {
			case monitor.DST_LOWER:
			case monitor.DST_HIGHER:
			}
		}
	}()

	log.Print("About to serve on port ", *argPort)
	addr := fmt.Sprintf(":%v", *argPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
