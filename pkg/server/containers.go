package server

import (
	"fmt"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) DeployJob(request *pb.DeployJobRequest, stream pb.Dokkup_DeployJobServer) error {
	s.Logger.Info("job received", zap.String("jobName", request.Name), zap.Int("containerCount", int(request.Count)), zap.String("containerImage", request.Container.Image))
	stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to pull image: %s", request.Container.Image)})

	s.Logger.Info("attempting to pull image", zap.String("image", request.Container.Image))
	err := s.Controller.ImagePull(request.Container.Image)
	if err != nil {
		s.Logger.Error("failed to pull image", zap.Error(err), zap.String("image", request.Container.Image))
		return fmt.Errorf("failed to pull image: %w", err)
	}

	shouldUpdate, err := s.Controller.ShouldUpdateContainers(request)
	if err != nil {
		s.Logger.Error("failed to check if an update should be run", zap.Error(err))
		return err
	}

	doesJobExist := s.Controller.DoesJobExist(request.Name)
	if doesJobExist && !shouldUpdate {
		s.Logger.Info("nothing to do, exiting...")
		return nil
	}

	createdContainers, err := s.Controller.CreateContainersFromRequest(request, stream)
	if err != nil {
		s.Logger.Error("failed to create containers", zap.Error(err))
		return err
	}

	if shouldUpdate {
		s.Logger.Info("updating the job...")
		jobContainers, err := s.Controller.GetContainersByJobName(request.Name)
		if err != nil {
			s.Logger.Error("failed to get containers by job name", zap.Error(err))
			return err
		}

		s.Logger.Info("stopping containers")
		err = s.Controller.StopContainers(jobContainers)
		if err != nil {
			s.Logger.Error("failed to stop containers", zap.Error(err))
			return err
		}
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

func (s *Server) StopJob(request *pb.StopJobRequest, stream pb.Dokkup_StopJobServer) error {
	s.Logger.Info("attempting to stop a job", zap.String("jobName", request.Name))

	s.Logger.Info("getting running containers")
	jobContainers, err := s.Controller.GetContainersByJobName(request.Name)
	if err != nil {
		s.Logger.Error("failed to get containers by job name", zap.Error(err))
		return err
	}

	s.Logger.Info("stopping containers")
	err = s.Controller.StopContainers(jobContainers)
	if err != nil {
		s.Logger.Error("failed to stop containers", zap.Error(err))
		return err
	}

	s.Logger.Info("deleting containers")
	err = s.Controller.DeleteContainers(jobContainers, stream)
	if err != nil {
		s.Logger.Error("failed to delete containers", zap.Error(err))
		return err
	}

	s.Logger.Info("job stopped successfully")
	stream.Send(&pb.StopJobResponse{Message: "Job stopped successfully"})
	return nil
}
