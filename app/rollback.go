package app

import (
	"context"
	"fmt"
	"github.com/XiovV/dokkup/controller"
	pb "github.com/XiovV/dokkup/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"time"
)

type Rollback struct {
	config     *Config
	controller controller.DockerController
	client pb.RollbackClient
}

func NewRollback(config *Config, dockerController controller.DockerController, conn *grpc.ClientConn) Rollback {
	return Rollback{
		config: config,
		controller: dockerController,
		client: pb.NewRollbackClient(conn),
	}
}

func (r *Rollback) rollbackContainer() (pb.Rollback_RollbackContainerClient, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)

	header := metadata.New(map[string]string{"authorization": r.config.APIKey})
	ctx = metadata.NewOutgoingContext(ctx, header)

	request := pb.RollbackRequest{Container: r.config.Container}
	stream, err := r.client.RollbackContainer(ctx, &request)
	if err != nil {
		log.Fatal(err)
	}

	return stream, cancel
}

func (r *Rollback) Run() {
	stream, cancel := r.rollbackContainer()
	defer cancel()

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
