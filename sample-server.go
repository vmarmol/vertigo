package main
 
import (
        "flag"
        "fmt"
        "log"
        "net/http"
	"time"
)
 
var argPort = flag.Int("port", 80, "port to listen")
var delay = flag.Int("delay", 250, "number of milliseconds to be busy.") 

func main() {
	startTime := time.Now()
	nanoDelay := int64(*delay * 1000 * 1000) // Milli to Nano
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqTime := time.Now()
		timeSince := time.Since(reqTime)
		for i := 0; timeSince.Nanoseconds() < nanoDelay; i++ {
			f := float64(i)
			n := float64(i+2)
			f = f / n
			n = n / f
			timeSince = time.Since(reqTime)
		}
		fmt.Fprintf(w, time.Since(startTime).String())
        })
 
        addr := fmt.Sprintf(":%v", *argPort)
        log.Fatal(http.ListenAndServe(addr, nil))
}
