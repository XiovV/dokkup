package main

import (
	"flag"
	"fmt"
	"github.com/XiovV/docker_control_cli/app"
	"os"
)


func main() {
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	statusNode := statusCmd.String("node", "", "Name of the node")
	statusGroup := statusCmd.String("group", "", "Name of the group")

	if len(os.Args) < 2 {
		fmt.Println("invalid number of arguments")
		os.Exit(1)
	}

	config := app.NewConfig("./config.json")

	app := app.New(config)

	switch os.Args[1] {
	case "status":
		statusCmd.Parse(os.Args[2:])
		app.HandleStatus(*statusNode, *statusGroup)
	default:
		fmt.Println("expected update or status command")
		os.Exit(1)
	}

	//config := app.NewConfig("./config.json")
	//
	//app := app.New(config)
	//
	//app.Start()
}