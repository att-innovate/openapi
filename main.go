package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	handler "openapi/comm"
	"openapi/conf"
	"openapi/core"
)

var openAPIPort = flag.Int("port", 8085, "Port used for the OpenAPI.")
var openAPIIP = flag.String("ip", "0.0.0.0", "IP address used for the OpenAPI.")
var mode = flag.String("mode", "lte", "Define the mode (lte|wifi).")

func main() {
	flag.Parse()
	config := conf.LoadConfiguration(*mode)

	var openAPI = core.OpenAPI{}
	openAPI.InitConfiguration(config)

	log.Printf("Invoke project 'OpenAPI' in version 1.0 in [%v] environment mode.", config.Env)

	listenOn := fmt.Sprintf("%s:%v", *openAPIIP, *openAPIPort)

	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/hello", handler.HelloHandler)

	http.HandleFunc("/statistics", handler.StatisticsHandler)
	http.HandleFunc("/streaming", func(w http.ResponseWriter, r *http.Request) { handler.StreamingModeHandler(w, r, config.Env) })
	http.HandleFunc("/latency", func(w http.ResponseWriter, r *http.Request) { handler.LatencyModeHandler(w, r, config.Env) })
	http.HandleFunc("/normal", func(w http.ResponseWriter, r *http.Request) { handler.NormalModeHandler(w, r, config.Env) })

	http.HandleFunc("/goodbye", handler.GoodbyHandler)

	log.Printf("Start server on %s.\n", listenOn)
	http.ListenAndServe(listenOn, nil)
}
