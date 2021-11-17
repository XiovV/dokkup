package app

import (
	"context"
	"fmt"
	"github.com/XiovV/dokkup/controller"
	pb "github.com/XiovV/dokkup/grpc"
	"io"
	"log"
	"time"
)

type Rollback struct {
	config     *Config
	controller controller.DockerController
}

func NewRollback(config *Config, dockerController controller.DockerController) Rollback {
	return Rollback{config: config, controller: dockerController}
}

func (r *Rollback) Run(client pb.RollbackClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)
	defer cancel()

	request := pb.RollbackRequest{Container: r.config.Container}
	stream, err := client.RollbackContainer(ctx, &request)
	if err != nil {
		log.Fatal(err)
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("err", err)
		}

		fmt.Println(response.GetMessage())
	}

}
