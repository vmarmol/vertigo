package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/user"
)

var argPort = flag.Int("port", 8080, "port")

func main() {
	flag.Parse()
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("unable to get current user: %v\n", err)
	}
	if currentUser.Username != "root" {
		log.Fatalf("must be root!")
	}
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
