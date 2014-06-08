package main

import (
	"log"
	"net/http"
)

func main() {
	taskManager, err := NewDockerTaskManager()
	if err != nil {
		log.Fatal(err)
	}
	rest := &restTaskManager{
		taskManager: taskManager,
	}
	http.Handle("/task", rest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
