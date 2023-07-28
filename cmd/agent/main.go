package main

import (
	"fmt"
	"log"

	"github.com/XiovV/dokkup/pkg/config"
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

  srv := server.Server{Config: config} 

  fmt.Println("server listening on port", config.Port)
  if err := srv.Serve(); err != nil {
    log.Fatal(err)
  }
}
