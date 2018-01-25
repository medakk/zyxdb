package raft

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ZyxdbConfig struct {
	Nodes []Node `yaml:"nodes"`
}

func (c *ZyxdbConfig) getNodeByName(name string) Node {
	for _, node := range c.Nodes {
		if node.Name == name {
			return node
		}
	}

	panic(fmt.Sprintf("No node with name %s in config", name))
}

func (c *ZyxdbConfig) nodeCount() int {
	return len(c.Nodes)
}

func loadConfig() (ZyxdbConfig, error) {
	c := ZyxdbConfig{}

	var f *os.File
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
		err := fmt.Errorf("no config file found. Tried: %v", configFilepaths)
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
