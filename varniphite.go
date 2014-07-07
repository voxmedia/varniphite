package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/marpaia/graphite-golang"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	graphiteHost string
	graphitePort int
	metricPath   string
	interval     int
)

func init() {
	flag.StringVar(&graphiteHost, "H", "localhost", "Hostname for graphite")
	flag.IntVar(&graphitePort, "p", 2003, "Port for graphite")
	flag.StringVar(&metricPath, "m", "varnish.stats", "Metric path")
	flag.IntVar(&interval, "i", 10, "Check stats each <i> seconds")
}

func main() {

	tenSecs := time.Duration(interval) * time.Second

	if len(os.Args) < 4 {
		flag.Usage()
	}
	flag.Parse()

	for {
		fmt.Printf("Running...")
		work()
		fmt.Printf("Done!\n")
		time.Sleep(tenSecs)
	}

}

func work() {
	// Read results from varnishstat
	out, err := exec.Command("varnishstat", "-j").Output()
	if err != nil {
		log.Fatalf("%v", err)
	}

	var data interface{}
	err = json.Unmarshal(out, &data)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Deal with this horrible json
	cleanData := make(map[string]float64)
	d := data.(map[string]interface{})

	for k, v := range d {
		switch i := v.(type) {
		case string:
			if strings.Contains(i, "timestamp") {
				continue
			}
		case map[string]interface{}:
			for sk, s := range i {
				switch si := s.(type) {
				case float64:
					if strings.Contains(sk, "value") {
						cleanData[k] = si
					}
				}
			}
		default:
			log.Println("nope")
		}
	}

	// Graphite connection and sending
	conn, err := graphite.NewGraphite(graphiteHost, graphitePort)
	if err != nil {
		log.Fatalf("%v", err)
	}

	re, err := regexp.Compile("\\(.*\\)")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for z, w := range cleanData {
		z := strings.ToLower(re.ReplaceAllString(z, ""))
		metric := fmt.Sprintf("%s.%s", metricPath, z)
		value := fmt.Sprintf("%f", w)
		conn.SimpleSend(metric, value)
	}
}
