package main

import (
	"flag"
	"github.com/XiovV/docker_control_cli/app"
	"github.com/XiovV/docker_control_cli/controller"
	"os"
)

func main() {
	var node, container, apiKey string

	actionCmd := flag.NewFlagSet("up", flag.ExitOnError)
	actionCmd.StringVar(&node, "node", "", "Node endpoint")
	actionCmd.StringVar(&container, "container", "", "Name of the container you wish to update")
	actionCmd.StringVar(&apiKey,"api-key", "", "Docker Control Agent API Key")
	keep := actionCmd.Bool("keep", false, "Keep the previous version of the container. Useful if you ever need to use the 'rollback' command")
	image := actionCmd.String("image", "", "The image you'd like the container to be updated to")

	rollbackCmd := flag.NewFlagSet("rollback", flag.ExitOnError)
	rollbackCmd.StringVar(&node, "node", "", "Node endpoint")
	rollbackCmd.StringVar(&container, "container", "", "Name of the container you wish to update")
	rollbackCmd.StringVar(&apiKey,"api-key", "", "Docker Control Agent API Key")

	switch os.Args[1] {
	case "up":
		actionCmd.Parse(os.Args[2:])
	case "rollback":
		rollbackCmd.Parse(os.Args[2:])
	}

	dockerController := controller.NewDockerController(node, apiKey)

	switch os.Args[1] {
	case "up":
		app := app.NewUpdate(node, container, *image, *keep, dockerController)
		app.Run()
	case "rollback":
		app := app.NewRollback(node, container, dockerController)
		app.Run()
	}

}
