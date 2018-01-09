package raft

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// How often to check whether or nodes are alive
	heartbeatDuration = 15 * time.Second
)

var (
	// tries reading configs in this order, with the first one that is present
	// being used
	configFilepaths = []string{"./zyxdb.yml", "/etc/zyxdb/zyxdb.yml"}
)

type RaftState struct {
	CurrentTerm int
	VotedFor    int
	Log         []string

	// Config where the list of nodes is stored
	Config ZyxdbConfig
}

type Node struct {
	Name     string `yaml: "name"`
	Location string `yaml: "location"`
}

type ZyxdbConfig struct {
	Nodes []Node `yaml: "nodes"`
}

func New() *RaftState {
	zyxdbConfig, err := loadConfig()
	if err != nil {
		panic(err.Error())
	}

	c := RaftState{
		Config: zyxdbConfig,
	}
	return &c
}

func (c *RaftState) Middleware(handler func(http.ResponseWriter, *http.Request, *RaftState)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}

func loadConfig() (ZyxdbConfig, error) {
	c := ZyxdbConfig{}

	var f *os.File = nil
	// Open the config file
	for _, filepath := range configFilepaths {
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			continue
		}

		configFile, err := os.Open(filepath)
		if err != nil {
			return c, err
		}

		f = configFile
		defer f.Close()
	}

	if f == nil {
		err := errors.New(fmt.Sprintf("No config file found. Tried: %v\n", configFilepaths))
		return c, err
	}

	// Read its contents
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return c, err
	}

	// Load the YAML into a struct
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
