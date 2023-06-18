package main

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// initialize .env variables
	config.InitVariables()

	// print out the port
	fmt.Println("[::] Starting server on the port: ", config.Port)

	r := router.Handler()

	// start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
