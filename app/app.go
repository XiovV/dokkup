package app

import (
	"fmt"
	"github.com/XiovV/docker_control_cli/controller"
	"os"
)

type App struct {
	action string
	node string
	container string
	image string
	controller controller.DockerController
}

func New(action, node, container, image string, controller controller.DockerController) *App {
	return &App{
		action:    action,
		node:      node,
		container: container,
		image:     image,
		controller: controller,
	}
}

func (a *App) Run() {
	fmt.Println("pulling image:", a.image)
	err := a.controller.PullImage(a.image)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("image pulled successfully")

	fmt.Printf("updating %s to %s\n", a.container, a.image)
	err = a.controller.UpdateContainer(a.container, a.image)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("container updated successfully")
}