package app

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const defaultConfigFilename = "dokkup.yaml"

type Config struct {
	APIKey       string `yaml:"api_key"`
	NodeLocation string `yaml:"node_location"`
	Container    string `yaml:"container"`
	Tag          string `yaml:"tag"`
	Image        string `yaml:"image"`
	Keep         bool   `yaml:"keep"`
}

func NewConfig() *Config {
	if len(os.Args) < 2 {
		fmt.Println("please provide a command (up or rollback)")
		os.Exit(1)
	}

	config := Config{}

	if len(os.Args) > 2 {
		config.parseCmdFlags()
		return &config
	}

	config.parseConfigFile(defaultConfigFilename)

	return &config
}

func (config *Config) parseConfigFile(filename string) {
	file, err := os.Open(filename)

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("A config file does not exist. Please create it or use command line flags.")
		os.Exit(1)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}
}

func (config *Config) parseCmdFlags() {
	switch os.Args[1] {
	case "up":
		actionCmd := flag.NewFlagSet("up", flag.ExitOnError)
		actionCmd.StringVar(&config.NodeLocation, "node", "", "Node endpoint")
		actionCmd.StringVar(&config.Container, "container", "", "Name of the container you wish to update")
		actionCmd.StringVar(&config.APIKey, "api-key", "", "Docker Control Agent API Key")
		actionCmd.BoolVar(&config.Keep, "keep", false, "Keep the previous version of the container. Useful if you ever need to use the 'rollback' command")
		actionCmd.StringVar(&config.Image, "image", "", "The image you'd like the container to be updated to")
		actionCmd.StringVar(&config.Tag, "tag", "", "Image tag you'd like the container to be updated to")

		if err := actionCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("flag parser error:", err)
		}

	case "rollback":
		rollbackCmd := flag.NewFlagSet("rollback", flag.ExitOnError)
		rollbackCmd.StringVar(&config.NodeLocation, "node", "", "Node endpoint")
		rollbackCmd.StringVar(&config.Container, "container", "", "Name of the container you wish to update")
		rollbackCmd.StringVar(&config.APIKey, "api-key", "", "Docker Control Agent API Key")

		if err := rollbackCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("flag parser error:", err)
		}
	}
}
