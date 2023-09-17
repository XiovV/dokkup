package server

import (
	"context"

	pb "github.com/XiovV/dokkup/pkg/grpc"
)

func (s *Server) GetNodeStatus(ctx context.Context, in *pb.GetNodeStatusRequest) (*pb.NodeStatus, error) {
	return nil, nil
}
