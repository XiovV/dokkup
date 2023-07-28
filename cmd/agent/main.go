package main

import (
	"fmt"
	"log"

	"github.com/XiovV/dokkup/pkg/config"
	"github.com/XiovV/dokkup/pkg/docker"
	"github.com/XiovV/dokkup/pkg/server"
)

func main() {
  key, err := config.CheckAPIKey()
  if err != nil {
    log.Fatal("could not check API key: ", err)
  }

  config, err := config.ReadAgentConfig()
  if err != nil {
    log.Fatal("could not read config: ", err)
  }

  config.APIKey = key

  controller, err := docker.NewController()
  if err != nil {
    log.Fatal("could not instantiate controller: ", err)
  } 

  srv := server.Server{Config: config, Controller: controller} 

  fmt.Println("server listening on port", config.Port)
  if err := srv.Serve(); err != nil {
    log.Fatal(err)
  }
}
