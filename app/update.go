package app

import (
	"fmt"
	"github.com/XiovV/dokkup/controller"
	"os"
	"strings"
)

type Update struct {
	config     *Config
	controller controller.DockerController
}

func NewUpdate(config *Config, controller controller.DockerController) *Update {
	return &Update{
		config:     config,
		controller: controller,
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

	if a.config.Tag != "" {
		image, err := a.controller.GetContainerImage(a.config.Container)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		imageParts := strings.Split(image, ":")

		if len(imageParts) != 2 || imageParts[0] == "" || imageParts[1] == "" {
			fmt.Println("agent returned an image with an invalid format")
			os.Exit(1)
		}

		a.config.Image = imageParts[0] + ":" + a.config.Tag
	}

	fmt.Println("pulling image:", a.config.Image)
	err := a.controller.PullImage(a.config.Image)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("image pulled successfully")

	fmt.Printf("updating %s to %s\n", a.config.Container, a.config.Image)
	err = a.controller.UpdateContainer(a.config.Container, a.config.Image, a.config.Keep)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("container updated successfully")
}
