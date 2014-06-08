package main
 
import (
	"flag"
	"fmt"
	"io/ioutil"
        "log"
        "net/http"
	"strconv"
	"time"
)

var argPort = flag.Int("port", 80, "port to listen")
var url = flag.String("url", "http://107.178.208.130:80", "Url to serve traffic to.")
var initQps = flag.Int("init_qps", 10, "Number of queries per second to initially send.") 

var nanoDelay = int64(1000 * 1000 * 1000) 
var uptime = "Initial Uptime (not updated.)"

func sendQueries() {
	for true {
		time.Sleep(time.Duration(nanoDelay))
		go SendOneQuery()
	}
}

func SendOneQuery() {
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

func main() {
	flag.Parse()
	nanoDelay = int64(1000 * 1000 * 1000 / *initQps)
	http.HandleFunc("/service/qps", func(w http.ResponseWriter, r *http.Request) {
		qps, err := strconv.Atoi(r.FormValue("qps"))		
		if err != nil {
                        log.Printf("%v", err)
			return
                }
		nanoDelay = int64(1000 * 1000 * 1000 / qps)
        })
	http.HandleFunc("/service/info", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, uptime)
        })
	go sendQueries()
	addr := fmt.Sprintf(":%v", *argPort)
        log.Fatal(http.ListenAndServe(addr, nil))
}
