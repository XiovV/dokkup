package server

import (
	"context"

	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) GetJobStatus(ctx context.Context, in *pb.Job) (*pb.JobStatus, error) {
	s.Logger.Info("retreiving job status")

	s.Logger.Debug("fetching running containers")
	runningContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Stopped: false})
	if err != nil {
		s.Logger.Error("could not get running containers", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("fetching all containers")
	totalContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		s.Logger.Error("could not get all containers", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("checking if an image already exists", zap.String("image", in.Container.Image))
	doesExist, err := s.Controller.ImageDoesExist(in.Container.Image)
	if err != nil {
		s.Logger.Error("failed to check if image exists", zap.Error(err), zap.String("image", in.Container.Image))
		return nil, err
	}

	s.Logger.Debug("does image exist", zap.Bool("doesExist", doesExist))
	shouldUpdate := !doesExist

	s.Logger.Debug("getting rollback containers")
	rollbackContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Rollback: true})
	if err != nil {
		s.Logger.Error("could not get rollback containers", zap.Error(err))
		return nil, err
	}

	canRollback := len(rollbackContainers) > 0

	response := &pb.JobStatus{
		TotalContainers:   int32(len(totalContainers)),
		RunningContainers: int32(len(runningContainers)),
		ShouldUpdate:      shouldUpdate,
		CanRollback:       canRollback,
	}

	s.Logger.Debug("should update job", zap.Bool("shouldUpdate", shouldUpdate))

	if len(totalContainers) == 0 {
		return response, nil
	}

	containerConfig, err := s.Controller.ContainerInspect(totalContainers[0].ID)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("creating temporary container")
	temporaryContainer, err := s.Controller.CreateTemporaryContainer(in)
	if err != nil {
		s.Logger.Error("failed to create temporary container", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("temporary container created successfully", zap.String("containerId", temporaryContainer))
	temporaryContainerConfig, err := s.Controller.ContainerInspect(temporaryContainer)
	if err != nil {
		s.Logger.Error("could not inspect temporary container", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("checking if configs are different")
	isDifferent := s.Controller.IsConfigDifferent(containerConfig, temporaryContainerConfig)

	s.Logger.Debug("is config different", zap.Bool("isDifferent", isDifferent))

	err = s.Controller.ContainerRemove(temporaryContainer)
	if err != nil {
		s.Logger.Error("failed to remove temporary container", zap.Error(err))
		return nil, err
	}

	response.ShouldUpdate = isDifferent
	s.Logger.Debug("should update job", zap.Bool("shouldUpdate", isDifferent))

	return response, nil
}
