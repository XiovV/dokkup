package server

import (
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

	for _, container := range currentContainers {
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

	for _, container := range rollbackContainers {
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

	for _, container := range currentContainers {
		err := s.Controller.ContainerRemove(container.ID)
		if err != nil {
			s.Logger.Error("could not remove container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	return nil
}
