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


  createdContainers := []string{}
	for i := 0; i < int(request.Count); i++ {
		uuid := uuid.New()
		containerName := fmt.Sprintf("%s-%s", request.Container.Name, uuid.String())
    
    stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Setting up container: %s", containerName)})

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

    resp, err := s.Controller.ContainerCreate(containerName, containerConfig, hostConfig)
		if err != nil {
			fmt.Println(err)
			return err
		}

    createdContainers = append(createdContainers, resp.ID)
    
    stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Container %s configured successfully", containerName)})
	}

  for _, container := range createdContainers {
    stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to start container: %s", container)})

    if err := s.Controller.ContainerStart(container); err != nil {
      return err
    }

    stream.Send(&pb.DeployJobResponse{Message: "Container started successfully"})
  }


  stream.Send(&pb.DeployJobResponse{Message: "Deployment successfull"})
	return nil
}
