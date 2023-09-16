package server

import (
	"fmt"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) RollbackJob(request *pb.RollbackJobRequest, stream pb.Dokkup_RollbackJobServer) error {
	s.Logger.Info("received a rollback request", zap.String("jobName", request.Name))
	currentContainers, err := s.Controller.GetContainersByJobName(request.Name)
	if err != nil {
		s.Logger.Error("could not get containers", zap.Error(err))
		return err
	}

	for i, container := range currentContainers {
		stream.Send(&pb.RollbackJobResponse{Message: fmt.Sprintf("Stopping container (%d/%d)", i+1, len(currentContainers))})
		err := s.Controller.ContainerStop(container.ID)
		if err != nil {
			s.Logger.Error("could not stop container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	rollbackContainers, err := s.Controller.GetRollbackContainers(request.Name)
	if err != nil {
		s.Logger.Error("could not get rollback containers", zap.Error(err))
		return err
	}

	for i, container := range rollbackContainers {
		stream.Send(&pb.RollbackJobResponse{Message: fmt.Sprintf("Starting container (%d/%d)", i+1, len(rollbackContainers))})
		err := s.Controller.ContainerStart(container.ID)
		if err != nil {
			s.Logger.Error("could not start rollback container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	// TODO: replace this with a method that will rename a singular container
	err = s.Controller.RemoveRollbackFromContainers(rollbackContainers)
	if err != nil {
		s.Logger.Error("could not remove remove -rollback from containers", zap.Error(err))
		return nil
	}

	for i, container := range currentContainers {
		stream.Send(&pb.RollbackJobResponse{Message: fmt.Sprintf("Removing container (%d/%d)", i+1, len(rollbackContainers))})
		err := s.Controller.ContainerRemove(container.ID)
		if err != nil {
			s.Logger.Error("could not remove container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	stream.Send(&pb.RollbackJobResponse{Message: "Rollback successful"})
	return nil
}
