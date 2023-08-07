package server

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) DeployContainer(ctx context.Context, in *pb.DeployContainerRequest) (*empty.Empty, error) {
	fmt.Println(in)

	err := s.Controller.ImagePull(in.ContainerImage)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to create container")
	}

	err = s.Controller.ContainerCreate(in.ContainerName, in.ContainerImage)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to create container")
	}

	return new(empty.Empty), nil
}
