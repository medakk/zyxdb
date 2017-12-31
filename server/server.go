package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.com/medakk/zyxdb/coordinator"
)

func ListenAndServe(name, host string, port int) {
	// Setup host:port to listen on
	listenOn := fmt.Sprintf("%s:%d", host, port)

	// Set up routing
	r := mux.NewRouter()
	r.HandleFunc("/retrieve/", RetrieveHandler)
	r.HandleFunc("/insert/", InsertHandler)
	r.HandleFunc("/ping/", PingHandler)

	// Set up Coordinator
	coord := coordinator.New()
	coordRouter := coord.Middleware(r)

	// Set up logging
	loggedRoute := handlers.LoggingHandler(os.Stderr, coordRouter)

	// Apply routing
	http.Handle("/", loggedRoute)

	log.Printf("Starting server %s on %s\n", name, listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
