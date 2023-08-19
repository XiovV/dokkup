package server

import (
	"fmt"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) DeployJob(request *pb.DeployJobRequest, stream pb.Dokkup_DeployJobServer) error {
	s.Logger.Info("job received", zap.String("jobName", request.Name), zap.Int("containerCount", int(request.Count)), zap.String("containerImage", request.Container.Image))

	doesJobExist, err := s.Controller.ContainerDoesExist(request.Name)
	if err != nil {
		s.Logger.Error("failed to check if deployment already exists", zap.Error(err))
		return err
	}

	if doesJobExist {
		s.Logger.Info("job already exists, run an update")
		return nil
	}

	s.Logger.Info("job doesn't exist, run a full deployment")

	stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to pull image: %s", request.Container.Image)})

	s.Logger.Info("attempting to pull image", zap.String("image", request.Container.Image))
	err = s.Controller.ImagePull(request.Container.Image)
	if err != nil {
		s.Logger.Error("failed to pull image", zap.Error(err), zap.String("image", request.Container.Image))
		return fmt.Errorf("failed to pull image: %w", err)
	}

	createdContainers, err := s.Controller.CreateContainersFromRequest(request, stream)
	if err != nil {
		s.Logger.Error("failed to create containers", zap.Error(err))
		return err
	}

	err = s.Controller.StartContainers(createdContainers, stream)
	if err != nil {
		s.Logger.Error("failed to start containers", zap.Error(err))
		return err
	}

	s.Logger.Info("job completed successfully")
	stream.Send(&pb.DeployJobResponse{Message: "Deployment successful"})
	return nil
}
