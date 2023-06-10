package main

import (
	"c2d-reports/internal/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Handler()

	// start the server
	fmt.Println("[::] Starting server on the port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
