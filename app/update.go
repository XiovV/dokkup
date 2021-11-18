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

func (u *Update) updateContainer() (pb.Updater_UpdateContainerClient, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)


	request := pb.UpdateRequest{Image: u.config.Tag, ContainerName: u.config.Container}
	stream, err := u.client.UpdateContainer(ctx, &request)
	if err != nil {
		log.Fatal(err)
	}

	return stream, cancel
}

func (u *Update) Run() {
	errors := u.ValidateFlags()
	if len(errors) != 0 {
		for _, error := range errors {
			fmt.Println(error)
		}

		os.Exit(1)
	}

	stream, cancel := u.updateContainer()
	defer cancel()

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
