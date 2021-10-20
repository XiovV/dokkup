package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Inventory struct {
	Groups []Group `yaml:"groups"`
}
type Group struct {
	Name  string   `yaml:"name"`
	Hosts []string `yaml:"hosts"`
}

func NewInventory(inventoryPath string) *Inventory {
	content, err := ioutil.ReadFile(inventoryPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var inventory Inventory
	err = yaml.Unmarshal(content, &inventory)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return &inventory
}

func (i *Inventory) FindGroupByName(name string) (Group, bool) {
	for _, group := range i.Groups {
		if group.Name == name {
			return group, true
		}
	}

	return Group{}, false
}