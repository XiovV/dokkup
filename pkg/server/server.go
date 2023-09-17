package server

import (
	"fmt"
	"net"

	"github.com/XiovV/dokkup/pkg/config"
	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/XiovV/dokkup/pkg/runner"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDokkupServer
	Config     *config.AgentConfig
	Controller *docker.Controller
	JobRunner  *runner.JobRunner
	Logger     *zap.Logger
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Config.Port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.authenticateUnary), grpc.StreamInterceptor(s.authenticateStream))
	pb.RegisterDokkupServer(grpcServer, s)

	return grpcServer.Serve(lis)
}
