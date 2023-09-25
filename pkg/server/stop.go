package server

import (
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"go.uber.org/zap"
)

func (s *Server) StopJob(request *pb.StopJobRequest, stream pb.Dokkup_StopJobServer) error {
	s.Logger.Info("stop job received", zap.String("jobName", request.Name))

	err := s.JobRunner.StopJob(request, stream)
	if err != nil {
		s.Logger.Error("stop job failed", zap.Error(err), zap.String("jobName", request.Name))
		return err
	}

	stream.Send(&pb.StopJobResponse{Message: "Job stopped successfully"})
	return nil
}
