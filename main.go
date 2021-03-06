package main

import (
	"fmt"
	"github.com/XiovV/dokkup/app"
	"github.com/XiovV/dokkup/controller"
	"os"
)

func main() {
	config := app.NewConfig()
	dockerController := controller.NewDockerController(config.NodeLocation, config.APIKey)

	switch os.Args[1] {
	case "up":
		app := app.NewUpdate(config, dockerController)
		app.Run()
	case "rollback":
		app := app.NewRollback(config, dockerController)
		app.Run()
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}
