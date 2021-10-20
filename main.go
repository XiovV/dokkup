package main

import (
	"flag"
	"fmt"
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

	err := actionCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dockerController := controller.NewDockerController(*node, *apiKey)
	app := app.New(os.Args[1], *node, *container, *image, dockerController)

	app.Run()
}
