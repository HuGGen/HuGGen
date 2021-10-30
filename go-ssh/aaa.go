package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/HuGGen/HuGGen/tree/master/go-ssh/util"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var configPath string
var servePort string

func main() {
	util.ParseFlags(&configPath, &servePort)

	http.HandleFunc("/", indexHandler)
	http.Handle("/metrics", promhttp.Handler())

	util.LogMsg(fmt.Sprintf("Listening on localhost:%s", servePort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", servePort), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	// Human readable navigation help.
	response := `<h1>ssh exporter</h1>
		<p><a href='/probe'>probe</a></p>
		<p><a href='/metrics'>metrics</a></p>`

	fmt.Fprintf(w, response)
}
