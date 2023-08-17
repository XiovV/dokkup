package server

import (
	"context"
	"fmt"
	"net"

	"github.com/XiovV/dokkup/pkg/config"
	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDokkupServer
	Config     *config.AgentConfig
	Controller *docker.Controller
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

func (s *Server) CheckAPIKey(ctx context.Context, in *pb.CheckAPIKeyRequest) (*empty.Empty, error) {
	// The interceptor already does the checking. This RPC is essentially used for 'pinging' nodes.
	return new(empty.Empty), nil
}
