package app

import (
	"context"
	"fmt"
	"github.com/XiovV/dokkup/controller"
	pb "github.com/XiovV/dokkup/grpc"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

type Update struct {
	config     *Config
	controller controller.DockerController
	client pb.UpdaterClient
}

func NewUpdate(config *Config, controller controller.DockerController, conn *grpc.ClientConn) *Update {
	return &Update{
		config:     config,
		controller: controller,
		client: pb.NewUpdaterClient(conn),
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

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)
	defer cancel()

	request := pb.UpdateRequest{Image: a.config.Tag, ContainerName: a.config.Container}
	stream, err := a.client.UpdateContainer(ctx, &request)
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
}
