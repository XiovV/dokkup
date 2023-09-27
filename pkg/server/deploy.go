package server

import (
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) DeployJob(request *pb.Job, stream pb.Dokkup_DeployJobServer) error {
	s.Logger.Info("job received", zap.String("jobName", request.Name), zap.Int("containerCount", int(request.Count)), zap.String("containerImage", request.Container.Image))

	s.Logger.Debug("checking if job already exists")
	doesJobExist := s.JobRunner.DoesJobExist(request.Name)

	s.Logger.Debug("does job already exist", zap.Bool("doesJobExist", doesJobExist))

	if doesJobExist {
		s.Logger.Debug("checking if the job should be updated")
		shouldUpdate, err := s.JobRunner.ShouldUpdateJob(request)
		if err != nil {
			s.Logger.Error("failed to check if an update should be run", zap.Error(err))
			return err
		}

		s.Logger.Debug("should update job", zap.Bool("shouldUpdate", shouldUpdate))

		if !shouldUpdate {
			s.Logger.Debug("nothing to do, exiting...")
			stream.Send(&pb.DeployJobResponse{Message: "Already up to date"})
			return nil
		}

		s.Logger.Info("updating job", zap.String("jobName", request.Name))
		return s.JobRunner.RunUpdate(request, stream)
	}

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
