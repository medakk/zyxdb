package main

import (
	"flag"

	"gitlab.com/medakk/zyxdb/server"
)

func main() {
	host := flag.String("host", "localhost", "host name to bind to")
	port := flag.Int("port", 5002, "port to listen on")
	flag.Parse()

	server.ListenAndServe(*host, *port)
}
