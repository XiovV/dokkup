package server

import (
	"context"
	"fmt"

	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) GetJobStatus(ctx context.Context, in *pb.GetJobStatusRequest) (*pb.JobStatus, error) {
	runningContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Stopped: false})
	if err != nil {
		s.Logger.Error("could not get running containers", zap.Error(err))
		return nil, err
	}

	totalContainers, err := s.Controller.GetContainers(in.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		s.Logger.Error("could not get all containers", zap.Error(err))
		return nil, err
	}

	fmt.Println("RUNNING CONTAINERS LEN", len(runningContainers))
	fmt.Println("TOTAL CONTAINERS LEN", len(totalContainers))

	response := &pb.JobStatus{
		TotalContainers:   int32(len(totalContainers)),
		RunningContainers: int32(len(runningContainers)),
		ShouldUpdate:      false,
	}

	return response, nil
}
