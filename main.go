package main

import (
	"flag"
	"github.com/XiovV/docker_control_cli/app"
	"github.com/XiovV/docker_control_cli/controller"
	"os"
)

func main() {
	actionCmd := flag.NewFlagSet("up", flag.ExitOnError)
	node := actionCmd.String("node", "", "Node endpoint")
	container := actionCmd.String("container", "", "Name of the container you wish to update")
	image := actionCmd.String("image", "", "The image you'd like the container to be updated to")
	apiKey := actionCmd.String("api-key", "", "Docker Control Agent API Key")
	keep := actionCmd.Bool("keep", false, "Keep the previous version of the container. Useful if you ever need to use the 'rollback' command")

	rollbackCmd := flag.NewFlagSet("rollback", flag.ExitOnError)
	rollbackNode := rollbackCmd.String("node", "", "Node endpoint")
	rollbackContainer := rollbackCmd.String("container", "", "Name of the container you wish to update")
	rollbackApiKey := rollbackCmd.String("api-key", "", "Docker Control Agent API Key")

	// TODO: clean this up
	switch os.Args[1] {
	case "up":
		actionCmd.Parse(os.Args[2:])
		dockerController := controller.NewDockerController(*node, *apiKey)
		app := app.NewUpdate(*node, *container, *image, *keep, dockerController)
		app.Run()
	case "rollback":
		rollbackCmd.Parse(os.Args[2:])
		dockerController := controller.NewDockerController(*rollbackNode, *rollbackApiKey)
		app := app.NewRollback(*rollbackNode, *rollbackContainer, dockerController)
		app.Run()
	}


}
