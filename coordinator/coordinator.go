package coordinator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	configFilepath = "./zyxdb.json"
)

type NodeConfig struct {
	Name     string `json: "name"`
	Location string `json: "location"`
}

type CoordinatorConfig struct {
	Nodes []NodeConfig `json: "nodes"`
}

type Coordinator struct {
	Config CoordinatorConfig
}

func (c *Coordinator) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%q\n", c.Config)
		h.ServeHTTP(w, r)
	})
}

func New() *Coordinator {
	config, err := loadConfig()
	if err != nil {
		panic(err.Error())
	}

	c := Coordinator{
		Config: config,
	}

	return &c
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
