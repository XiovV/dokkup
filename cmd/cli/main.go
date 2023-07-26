package main

import (
	"log"

	"github.com/XiovV/dokkup/cmd/cli/app"
	"github.com/XiovV/dokkup/config"
)

func main() {
  inventory, err := config.ReadInventory("../../inventory.yaml")
  if err != nil {
    log.Fatal(err)
  }

  app := app.NewApp(inventory)

  app.Run()
}


