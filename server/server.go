package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/medakk/zyxdb/raft"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ListenAndServe(name, host string, port int) {
	// Setup host:port to listen on
	listenOn := fmt.Sprintf("%s:%d", host, port)

	// Set up raft
	raft := raft.New(name)

	// Set up routing
	r := mux.NewRouter()
	r.HandleFunc("/ping/", PingHandler).Methods("GET")
	r.HandleFunc("/append-entries/", raft.Middleware(AppendEntriesHandler)).Methods("POST")
	r.HandleFunc("/request-vote/", raft.Middleware(RequestVoteHandler)).Methods("POST")
	r.HandleFunc("/insert/", raft.Middleware(InsertToLogHandler)).Methods("POST")
	r.HandleFunc("/debug-get-entries/", raft.Middleware(DebugGetEntriesHandler)).Methods("POST")

	// Set up logging
	loggedRoute := handlers.LoggingHandler(os.Stderr, r)

	// Apply routing
	http.Handle("/", loggedRoute)

	log.Printf("Starting server %s on %s\n", name, listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
