package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config []struct {
	Name       string   `yaml:"name"`
	Hosts      string   `yaml:"hosts"`
	Image      string   `yaml:"image"`
	Containers []string `yaml:"containers"`
}

func NewConfig(configPath string) *Config {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return &config
}