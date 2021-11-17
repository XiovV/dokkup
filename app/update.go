package app

import (
	"context"
	"fmt"
	"github.com/XiovV/dokkup/controller"
	pb "github.com/XiovV/dokkup/grpc"
	"io"
	"log"
	"os"
	"time"
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

func (a *Update) Run(client pb.UpdaterClient) {
	errors := a.ValidateFlags()
	if len(errors) != 0 {
		for _, error := range errors {
			fmt.Println(error)
		}

		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)
	defer cancel()

	request := pb.UpdateRequest{Image: a.config.Tag, ContainerName: a.config.Container}
	stream, err := client.UpdateContainer(ctx, &request)
	if err != nil {
		fmt.Println(err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response.GetMessage())
	}

	//err := a.controller.PullImage(a.config.Image)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//fmt.Println("image pulled successfully")
	//
	//fmt.Printf("updating %s to %s\n", a.config.Container, a.config.Image)
	//err = a.controller.UpdateContainer(a.config.Container, a.config.Image, a.config.Keep)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//fmt.Println("container updated successfully")
}
