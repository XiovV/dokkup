package app

import (
	"fmt"
	"github.com/XiovV/docker_control_cli/controller"
	"os"
)

type Update struct {
	node string
	container string
	image string
	keep bool
	controller controller.DockerController
}

func NewUpdate(node, container, image string, keep bool, controller controller.DockerController) *Update {
	return &Update{
		node:      node,
		container: container,
		image:     image,
		controller: controller,
		keep: keep,
	}
}

func (a *Update) Run() {
	fmt.Println("pulling image:", a.image)
	err := a.controller.PullImage(a.image)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("image pulled successfully")

	fmt.Printf("updating %s to %s\n", a.container, a.image)
	err = a.controller.UpdateContainer(a.container, a.image, a.keep)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("container updated successfully")
}