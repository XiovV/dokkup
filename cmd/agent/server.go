package main

import (
	"fmt"
	"net"

	"github.com/XiovV/dokkup/config"
	pb "github.com/XiovV/dokkup/grpc"
	"google.golang.org/grpc"
)

type Server struct {
  pb.UnimplementedDokkupServer
  Config *config.AgentConfig
}

func (s *Server) Serve() error {
  lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Config.Port))
  if err != nil {
    return err
  }

  grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.authenticateUnary))
  pb.RegisterDokkupServer(grpcServer, s)

  return grpcServer.Serve(lis)
}
