package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var argPort = flag.Int("port", 8080, "port")

func main() {
	flag.Parse()
	taskManager, err := NewDockerTaskManager()
	if err != nil {
		log.Fatal(err)
	}
	rest := &restTaskManager{
		taskManager: taskManager,
	}
	export := &taskExport{
		taskManager: taskManager,
	}
	http.Handle("/task", rest)
	http.Handle("/export/", export)
	addr := fmt.Sprintf(":%v", *argPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
