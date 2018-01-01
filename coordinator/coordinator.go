package coordinator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	// TODO: Use YAML or something more suited
	configFilepath = "/etc/zyxdb/zyxdb.json"

	// How often to check whether or nodes are alive
	heartbeatDuration = 15 * time.Second
)

type Node struct {
	Name     string `json: "name"`
	Location string `json: "location"`
}

type CoordinatorConfig struct {
	Nodes []Node `json: "nodes"`
}

type Coordinator struct {
	SelfNode Node
	Config   CoordinatorConfig

	// This stores whether or not all the nodes are alive
	allNodesOk     bool
	allNodesOkLock sync.RWMutex
}

// Initialize a new coordinator
func New(name, location string) *Coordinator {
	config, err := loadConfig()
	if err != nil {
		panic(err.Error())
	}

	c := Coordinator{
		SelfNode: Node{
			Name:     name,
			Location: location,
		},
		Config: config,

		// assume that everything is bad until proved otherwise
		allNodesOk:     false,
		allNodesOkLock: sync.RWMutex{},
	}

	go c.Heartbeat()

	return &c
}

func (c *Coordinator) Heartbeat() {
	for {
		allNodesOk := true

		// TODO: check health of all nodes concurrently
		// TODO: I am checking my own health. Do I need to do this?
		for _, node := range c.Config.Nodes {
			// TODO: Use a proper way to join URLs
			resp, err := http.Get(fmt.Sprintf("http://%s/ping/", node.Location))
			if err != nil {
				log.Printf("error with node %s, at %s: %s\n", node.Name, node.Location, err.Error())
				allNodesOk = false
				continue
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("error with node %s, at %s: %s\n", node.Name, node.Location)
				allNodesOk = false
				continue
			}
		}

		c.allNodesOkLock.Lock()
		c.allNodesOk = allNodesOk
		c.allNodesOkLock.Unlock()

		time.Sleep(heartbeatDuration)
	}
}

func (c *Coordinator) AllNodesOk() bool {
	c.allNodesOkLock.RLock()
	defer c.allNodesOkLock.RUnlock()

	return c.allNodesOk
}

func (c *Coordinator) Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the state of the heartbeats
		if !c.AllNodesOk() {
			w.WriteHeader(http.StatusServiceUnavailable)
			msg, _ := json.Marshal(map[string]string{
				"status": "heartbeat failed",
			})
			w.Write(msg)
			return
		}

		// Check which node should handle this
		vars := mux.Vars(r)
		key := vars["key"]
		nodeCount := len(c.Config.Nodes)
		targetNode := c.Config.Nodes[int(hash(key))%nodeCount]

		// TODO: The assertion of unique names must be done before this logic
		// can be used reliably
		if targetNode.Name != c.SelfNode.Name {
			targetUrl, _ := url.Parse(fmt.Sprintf("http://%s/", targetNode.Location))
			revProxy := httputil.NewSingleHostReverseProxy(targetUrl)
			revProxy.ServeHTTP(w, r)
			return
		}

		// This node should handle this. Fall through
		h.ServeHTTP(w, r)
	}
}

func loadConfig() (CoordinatorConfig, error) {
	c := CoordinatorConfig{}

	// Open the config file
	f, err := os.Open(configFilepath)
	if err != nil {
		return c, err
	}
	defer f.Close()

	// Read its contents
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return c, err
	}

	// Load the JSON
	err = json.Unmarshal(content, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func hash(s string) uint32 {
	var val uint32 = 31
	for pos, char := range s {
		val = val * (1 + (uint32(pos) * uint32(char)))
	}

	return val
}
