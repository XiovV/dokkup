package server

import (
	"errors"
	"fmt"

	pb "github.com/XiovV/dokkup/pkg/grpc"
)

func (s *Server) DeployJob(request *pb.DeployJobRequest, stream pb.Dokkup_DeployJobServer) error {
	stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to pull image: %s", request.Container.Image)})

	err := s.Controller.ImagePull(request.Container.Image)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to pull image")
	}

	createdContainers, err := s.Controller.CreateContainersFromRequest(request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.Controller.StartContainers(createdContainers)
	if err != nil {
		fmt.Println(err)
		return err
	}

	stream.Send(&pb.DeployJobResponse{Message: "Deployment successfull"})
	return nil
}
