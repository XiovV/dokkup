package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Groups []Groups `json:"groups"`
}

type Node struct {
	Location string `json:"location"`
	NodeName string `json:"node_name"`
}

type Groups struct {
	Group      string      `json:"group"`
	Containers []string    `json:"containers"`
	Image string `json:"image"`
	Nodes []Node `json:"nodes"`
}

func NewConfig(configPath string) *Config {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return &config
}

func (c *Config) FindNodeByName(nodeName string) (Node, bool) {
	for _, group := range c.Groups {
		for _, node := range group.Nodes {
			if node.NodeName == nodeName {
				return node, true
			}
		}
	}

	return Node{}, false
}

func (c *Config) FindGroupByName(groupName string) Groups {
	for _, group := range c.Groups {
		if group.Group == groupName {
			return group
		}
	}

	return Groups{}
}