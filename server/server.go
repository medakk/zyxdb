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

	// Set up Coordinator
	coord := coordinator.New(name, listenOn)

	// Set up routing
	r := mux.NewRouter()
	r.HandleFunc("/ping/", PingHandler).Methods("GET")
	r.HandleFunc("/retrieve/{key}/", coord.Middleware(RetrieveHandler)).Methods("POST")
	r.HandleFunc("/insert/{key}/", coord.Middleware(InsertHandler)).Methods("POST")

	// Set up logging
	loggedRoute := handlers.LoggingHandler(os.Stderr, r)

	// Apply routing
	http.Handle("/", loggedRoute)

	log.Printf("Starting server %s on %s\n", name, listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
