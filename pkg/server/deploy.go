package server

import (
	"fmt"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) DeployJob(request *pb.DeployJobRequest, stream pb.Dokkup_DeployJobServer) error {
	s.Logger.Info("job received", zap.String("jobName", request.Name), zap.Int("containerCount", int(request.Count)), zap.String("containerImage", request.Container.Image))

	stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to pull image: %s", request.Container.Image)})

	s.Logger.Debug("attempint to pull image", zap.String("image", request.Container.Image))
	err := s.Controller.ImagePull(request.Container.Image)
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	s.Logger.Debug("creating temporary container")
	temporaryContainer, err := s.Controller.CreateTemporaryContainer(request)
	if err != nil {
		s.Logger.Error("failed to create temporary container", zap.Error(err))
		return err
	}

	temporaryContainerConfig, err := s.Controller.ContainerInspect(temporaryContainer)
	if err != nil {
		s.Logger.Error("could not inspect temporary container", zap.Error(err))
		return err
	}

	s.Logger.Debug("checking if the job should be updated")
	shouldUpdate, err := s.JobRunner.ShouldUpdateJob(request.Name, temporaryContainerConfig)
	if err != nil {
		s.Logger.Error("failed to check if an update should be run", zap.Error(err))
		return err
	}

	s.Logger.Debug("removing the temporary container")
	err = s.Controller.ContainerRemove(temporaryContainer)
	if err != nil {
		s.Logger.Error("failed to remove the temporary container", zap.Error(err))
		return err
	}

	if shouldUpdate {
		s.Logger.Info("updating job", zap.String("jobName", request.Name))
		return s.JobRunner.RunUpdate(request, stream)
	}

	s.Logger.Debug("checking if job already exists")
	doesJobExist := s.JobRunner.DoesJobExist(request.Name)

	if doesJobExist {
		s.Logger.Debug("nothing to do, exiting...")
		return nil
	}

	s.Logger.Info("deploying a new job")
	err = s.JobRunner.RunDeployment(stream, request)
	if err != nil {
		s.Logger.Error("could not deploy job", zap.Error(err))
		return err
	}

	s.Logger.Info("job completed successfully")
	stream.Send(&pb.DeployJobResponse{Message: "Deployment successful"})
	return nil
}
