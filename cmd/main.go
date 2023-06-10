package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// start the server
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	}).Methods("GET")

	fmt.Println("[::] Starting server on the port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
