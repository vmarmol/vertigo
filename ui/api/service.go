package api

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var url = flag.String("url", "http://107.178.208.130:80", "Url to serve traffic to.")
var initQps = flag.Int("init_qps", 1, "Number of queries per second to initially send.")

var nsInSecond = int64(1000 * 1000 * 1000)
var nanoDelay = nsInSecond
var uptime = "Initial Uptime (not updated.)"

func sendQueries() {
	for true {
		time.Sleep(time.Duration(nanoDelay))
		go sendOneQuery()
	}
}

func sendOneQuery() {
	resp, err := http.Get(*url)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	uptime = string(body)
	log.Printf("New Uptime: %s", uptime)
}

func RegisterServiceHandlers() {
	nanoDelay = nsInSecond / int64(*initQps)
	http.HandleFunc("/api/qps", func(w http.ResponseWriter, r *http.Request) {
		qps, err := strconv.Atoi(r.FormValue("qps"))
		if err != nil {
			log.Printf("%v", err)
			return
		}
		nanoDelay = nsInSecond / int64(qps)
	})
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, uptime)
	})
	go sendQueries()
}
