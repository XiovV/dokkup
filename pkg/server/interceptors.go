package server

import (
	"context"

	"go.uber.org/zap"
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
		s.Logger.Warn("api key not provided")
		return status.Error(codes.Unauthenticated, "api key not provided")
	}

	apiKey := apiKeyHeader[0]

	if len(apiKey) > API_KEY_LENGHT {
		s.Logger.Warn("api key is too long", zap.Int("expectedLength", API_KEY_LENGHT), zap.Int("got", len(apiKey)))
		return status.Error(codes.Unauthenticated, "api key is invalid")
	}

	if len(apiKey) < API_KEY_LENGHT {
		s.Logger.Warn("api key is too short", zap.Int("expectedLength", API_KEY_LENGHT), zap.Int("got", len(apiKey)))
		return status.Error(codes.Unauthenticated, "api key is invalid")
	}

	key := s.Config.APIKey

	err := bcrypt.CompareHashAndPassword([]byte(key), []byte(apiKey))
	if err != nil {
		s.Logger.Warn("api key is incorrect")
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

func (s *Server) authenticateStream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := s.validateAuthHeader(stream.Context())
	if err != nil {
		return err
	}

	return handler(srv, stream)
}
