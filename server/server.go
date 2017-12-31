package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ListenAndServe(host string, port int) {
	// Setup host:port to listen on
	listenOn := fmt.Sprintf("%s:%d", host, port)

	// Set up routing
	r := mux.NewRouter()
	r.HandleFunc("/retrieve/", RetrieveHandler)
	r.HandleFunc("/insert/", InsertHandler)

	// Set up logging
	loggedRoute := handlers.LoggingHandler(os.Stderr, r)

	// Apply routing
	http.Handle("/", loggedRoute)

	log.Printf("Starting server on %s\n", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
