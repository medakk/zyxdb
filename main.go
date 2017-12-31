package main

import (
	"flag"
	"math/rand"
	"time"

	"gitlab.com/medakk/zyxdb/server"
)

func genRandomName() string {
	const chars = "abcdef0123456789"
	b := make([]byte, 12)
	b[0], b[1], b[2], b[3] = 'z', 'y', 'x', '-'
	for i := 4; i < 12; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

func main() {
	// Setup command line flags
	name := flag.String("name", "<random>", "name for this zyx instance")
	host := flag.String("host", "localhost", "host name to bind to")
	port := flag.Int("port", 5002, "port to listen on")
	flag.Parse()

	// initialize random number seed
	rand.Seed(time.Now().UTC().UnixNano())

	// generate random number if required
	if *name == "<random>" {
		*name = genRandomName()
	}

	server.ListenAndServe(*name, *host, *port)
}
