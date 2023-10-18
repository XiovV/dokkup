package server

import (
	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) DeployJob(request *pb.Job, stream pb.Dokkup_DeployJobServer) error {
	s.Logger.Info("job received", zap.String("jobName", request.Name), zap.Int("containerCount", int(request.Count)), zap.String("containerImage", request.Container.Image))

	s.Logger.Debug("checking if job already exists")
	doesJobExist := s.JobRunner.DoesJobExist(request.Name)

	s.Logger.Debug("does job already exist", zap.Bool("doesJobExist", doesJobExist))
	if !doesJobExist {
		s.Logger.Info("deploying a new job")
		err := s.JobRunner.RunDeployment(stream, request)
		if err != nil {
			s.Logger.Error("could not deploy job", zap.Error(err))
			return err
		}

		s.Logger.Info("job completed successfully")
		stream.Send(&pb.DeployJobResponse{Message: "Deployment successful"})
		return nil
	}

	currentContainers, err := s.Controller.GetContainers(request.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		s.Logger.Error("could not get current containers", zap.Error(err))
		return err
	}

	if int(request.Count) > len(currentContainers) {
		s.Logger.Debug("upscaling job", zap.Int("currentCount", len(currentContainers)), zap.Int("targetCount", int(request.Count)))
		count := int(request.Count) - len(currentContainers)
		return s.JobRunner.UpscaleJob(count, request, stream)
	}

	s.Logger.Debug("checking if the job should be updated")
	shouldUpdate, err := s.JobRunner.ShouldUpdateJob(request, currentContainers)
	if err != nil {
		s.Logger.Error("failed to check if an update should be run", zap.Error(err))
		return err
	}

	s.Logger.Debug("should update job", zap.Bool("shouldUpdate", shouldUpdate))

	if shouldUpdate {
		s.Logger.Info("updating job", zap.String("jobName", request.Name))
		return s.JobRunner.RunUpdate(request, stream)
	}

	stoppedContainers, err := s.Controller.GetStoppedContainers(request.Name)
	if err != nil {
		s.Logger.Error("could not get stopped containers", zap.Error(err))
		return err
	}

	if len(stoppedContainers) == 0 {
		s.Logger.Debug("nothing to do, exiting...")
		stream.Send(&pb.DeployJobResponse{Message: "Already up to date"})
		return nil
	}

	err = s.JobRunner.StartContainers(stoppedContainers, stream)
	if err != nil {
		s.Logger.Error("could not start container", zap.Error(err))
		return err
	}

	stream.Send(&pb.DeployJobResponse{Message: "Containers started successfully"})

	return nil
}
