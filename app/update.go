package app

import (
	"fmt"
	"github.com/XiovV/dokkup/controller"
	"os"
	"strings"
)

type Update struct {
	node       string
	container  string
	image      string
	tag        string
	keep       bool
	controller controller.DockerController
}

func NewUpdate(node, container, image, tag string, keep bool, controller controller.DockerController) *Update {
	return &Update{
		node:       node,
		container:  container,
		image:      image,
		tag:        tag,
		controller: controller,
		keep:       keep,
	}
}

func (a *Update) Run() {
	errors := a.ValidateFlags()
	if len(errors) != 0 {
		for _, error := range errors {
			fmt.Println(error)
		}

		os.Exit(1)
	}

	if a.tag != "" {
		image, err := a.controller.GetContainerImage(a.container)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		imageParts := strings.Split(image, ":")

		if len(imageParts) != 2 || imageParts[0] == "" || imageParts[1] == "" {
			fmt.Println("agent returned an image with an invalid format")
			os.Exit(1)
		}

		a.image = imageParts[0] + ":" + a.tag
	}

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
