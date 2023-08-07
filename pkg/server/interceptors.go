package server

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Server) validateAuthHeader(ctx context.Context) error {
	header, _ := metadata.FromIncomingContext(ctx)

	apiKeyHeader := header["authorization"]

	if len(apiKeyHeader) < 1 {
		fmt.Println("api key not provided")
		return status.Error(codes.Unauthenticated, "api key not provided")
	}

	apiKey := apiKeyHeader[0]

	if len(apiKey) > API_KEY_LENGHT {
		fmt.Println("api key too long")
		return status.Error(codes.Unauthenticated, "api key is invalid")
	}

	if len(apiKey) < API_KEY_LENGHT {
		fmt.Println("api key too short")
		return status.Error(codes.Unauthenticated, "api key is invalid")
	}

	key := s.Config.APIKey

	err := bcrypt.CompareHashAndPassword([]byte(key), []byte(apiKey))
	if err != nil {
		fmt.Println("api key incorrect")
		return status.Error(codes.Unauthenticated, "api key is invalid")
	}

	return nil
}

func (s *Server) authenticateUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := s.validateAuthHeader(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
