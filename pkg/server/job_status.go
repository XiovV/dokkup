package server

import (
	"context"

	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/XiovV/dokkup/pkg/version"
	"go.uber.org/zap"
)

func (s *Server) GetJobStatus(ctx context.Context, in *pb.Job) (*pb.JobStatus, error) {
	s.Logger.Info("retreiving job status")

	s.Logger.Debug("fetching all containers")
	totalContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		s.Logger.Error("could not get all containers", zap.Error(err))
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

	defer s.Controller.ContainerRemove(temporaryContainer)

	newVersionHash := version.Hash(temporaryContainerConfig)

	if len(totalContainers) == 0 {
		response := &pb.JobStatus{
			TotalContainers:   0,
			RunningContainers: 0,
			ShouldUpdate:      true,
			CanRollback:       false,
			CurrentVersion:    "",
			NewVersion:        newVersionHash,
		}

		return response, nil
	}

	s.Logger.Debug("fetching running containers")
	runningContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Stopped: false})
	if err != nil {
		s.Logger.Error("could not get running containers", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("getting rollback containers")
	rollbackContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Rollback: true})
	if err != nil {
		s.Logger.Error("could not get rollback containers", zap.Error(err))
		return nil, err
	}

	containerConfig, err := s.Controller.ContainerInspect(totalContainers[0].ID)
	if err != nil {
		return nil, err
	}

	currentVersionHash := version.Hash(containerConfig)

	s.Logger.Debug("checking if configs are different")
	isDifferent := s.Controller.IsConfigDifferent(containerConfig, temporaryContainerConfig)

	s.Logger.Debug("is config different", zap.Bool("isDifferent", isDifferent))

	s.Logger.Debug("should update job", zap.Bool("shouldUpdate", isDifferent))

	response := &pb.JobStatus{
		TotalContainers:   int32(len(totalContainers)),
		RunningContainers: int32(len(runningContainers)),
		ShouldUpdate:      isDifferent,
		CanRollback:       len(rollbackContainers) != 0,
		CurrentVersion:    currentVersionHash,
		NewVersion:        newVersionHash,
	}

	return response, nil
}
