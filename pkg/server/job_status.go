package server

import (
	"context"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/XiovV/dokkup/pkg/version"
	"go.uber.org/zap"
)

func (s *Server) GetJobStatus(ctx context.Context, in *pb.Job) (*pb.JobStatus, error) {
	s.Logger.Info("retreiving job status")

	jobContainers, err := s.JobRunner.GetJobStatus(in)
	if err != nil {
		s.Logger.Error("could not get job status", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("creating temporary container")
	temporaryContainer, err := s.Controller.CreateTemporaryContainer(in)
	if err != nil {
		s.Logger.Error("failed to create temporary container", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("temporary container created successfully", zap.String("containerId", temporaryContainer.ID))

	newVersionHash := version.Hash(temporaryContainer)

	if jobContainers.TotalContainers == 0 {
		response := &pb.JobStatus{
			ShouldUpdate: true,
			NewVersion:   newVersionHash,
		}

		return response, nil
	}

	isDifferent := s.Controller.IsConfigDifferent(jobContainers.RunningContainerConfig, temporaryContainer)

	response := &pb.JobStatus{
		TotalContainers:   int32(jobContainers.TotalContainers),
		RunningContainers: int32(jobContainers.RunningContainers),
		ShouldUpdate:      isDifferent,
		CurrentVersion:    version.Hash(jobContainers.RunningContainerConfig),
		NewVersion:        newVersionHash,
	}

	if jobContainers.RollbackContainers == 0 {
		response.CanRollback = false
		return response, nil
	}

	response.CanRollback = true
	response.OldVersion = version.Hash(jobContainers.RollbackContainerConfig)

	return response, nil
}
