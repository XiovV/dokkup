package server

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/XiovV/dokkup/pkg/version"
	"go.uber.org/zap"
)

func (s *Server) GetJobStatus(ctx context.Context, job *pb.Job) (*pb.JobStatus, error) {
	s.Logger.Info("retreiving job status")

	jobStatus, err := s.JobRunner.GetJobStatus(job)
	if err != nil {
		s.Logger.Error("could not get job status", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("creating temporary container")
	temporaryContainer, err := s.Controller.CreateTemporaryContainer(job)
	if err != nil {
		s.Logger.Error("failed to create temporary container", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("temporary container created successfully", zap.String("containerId", temporaryContainer.ID))

	newVersionHash := version.Hash(temporaryContainer)

	if len(jobStatus.TotalContainers) == 0 {
		response := &pb.JobStatus{
			ShouldUpdate: true,
			NewVersion:   newVersionHash,
		}

		return response, nil
	}

	isDifferent := s.Controller.IsConfigDifferent(jobStatus.RunningContainerConfig, temporaryContainer)

	containers := []*pb.ContainerInfo{}
	for _, container := range jobStatus.TotalContainers {
		ports := []string{}
		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprintf("%s:%d->%d/%s", port.IP, port.PublicPort, port.PrivatePort, port.Type))
		}

		containers = append(containers, &pb.ContainerInfo{
			Id:     container.ID,
			Name:   container.Names[0][1:],
			Status: container.Status,
			Ports:  strings.Join(ports, ""),
		})
	}

	if len(jobStatus.TotalContainers) != int(job.Count) {
		isDifferent = true
	}

	response := &pb.JobStatus{
		TotalContainers:   int32(len(jobStatus.TotalContainers)),
		RunningContainers: int32(len(jobStatus.RunningContainers)),
		ShouldUpdate:      isDifferent,
		CurrentVersion:    version.Hash(jobStatus.RunningContainerConfig),
		NewVersion:        newVersionHash,
		Containers:        containers,
		Image:             jobStatus.Image,
	}

	if len(jobStatus.RollbackContainers) == 0 {
		response.CanRollback = false
		return response, nil
	}

	response.CanRollback = true
	response.OldVersion = version.Hash(jobStatus.RollbackContainerConfig)

	return response, nil
}
