package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Inventory struct {
  Nodes []Node `yaml:"nodes"`
  Groups []Group `yaml:"groups"`
}

type Node struct {
  Name string `yaml:"name"`
  Location string `yaml:"location"`
  Key string `yaml:"key"`
}

type Group struct {
  Name string `yaml:"name"`
  Nodes []string `yaml:"nodes"`
}

func ReadInventory(filename string) (*Inventory, error) {
  content, err := ioutil.ReadFile(filename) 
  if err != nil {
    return nil, err
  }

  inventory := Inventory{} 
  err = yaml.Unmarshal(content, &inventory)
  if err != nil {
    return nil, err
  }

  err = inventory.validateGroups()
  if err != nil {
    return nil, err
  }
  
  return &inventory, nil
}

func (i *Inventory) validateGroups() error {
  for _, group := range i.Groups {
    for _, node := range group.Nodes {
      exists := i.doesNodeExist(node)

      if !exists {
        return fmt.Errorf("the node %s defined in group %s is not defined", node, group.Name)
      }
    }
  }

  return nil
}

func (i *Inventory) doesNodeExist(node string) bool {
  for _, n := range i.Nodes {
    if n.Name == node {
      return true
    } 
  } 

  return false
}
