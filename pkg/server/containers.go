package server

import (
	"errors"
	"fmt"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

func (s *Server) DeployJob(request *pb.DeployJobRequest, stream pb.Dokkup_DeployJobServer) error {
	stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to pull image: %s", request.Container.Image)})
	err := s.Controller.ImagePull(request.Container.Image)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to pull image")
	}

	for i := 0; i < int(request.Count); i++ {

		containerConfig := &container.Config{
			Image:        request.Container.Image,
			ExposedPorts: nat.PortSet{},
		}

		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{},
		}

		for _, port := range request.Container.Ports {
			internalPort := fmt.Sprintf("%s/tcp", port.In)
			containerConfig.ExposedPorts[nat.Port(internalPort)] = struct{}{}

			hostConfig.PortBindings[nat.Port(internalPort)] = []nat.PortBinding{{HostIP: "0.0.0.0"}}
		}

		uuidv4 := uuid.New()
		containerName := fmt.Sprintf("%s-%s", request.Container.Name, uuidv4.String())

		err = s.Controller.ContainerCreate(containerName, containerConfig, hostConfig)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("successfully created container")

	}

	return nil
}
