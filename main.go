package main

import (
	"flag"
	"fmt"
	"github.com/XiovV/docker_control_cli/app"
	"github.com/XiovV/docker_control_cli/services"
	"os"
)


func main() {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	groupUpdate := updateCmd.String("group", "", "Name of the node")

	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	statusNode := statusCmd.String("node", "", "Name of the node")
	statusGroup := statusCmd.String("group", "", "Name of the group")

	if len(os.Args) < 2 {
		fmt.Println("invalid number of arguments")
		os.Exit(1)
	}

	config := app.NewConfig("./config.json")
	dockerService := services.NewDockerController()

	app := app.New(config, dockerService)

	switch os.Args[1] {
	case "status":
		if err := statusCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("error parsing status command:", err)
			os.Exit(1)
		}
		app.HandleStatus(*statusNode, *statusGroup)
	case "update":
		if err := updateCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("error parsing update command:", err)
			os.Exit(1)
		}
		app.HandleUpdate(*groupUpdate)
	default:
		fmt.Println("expected update or status command")
		os.Exit(1)
	}
}