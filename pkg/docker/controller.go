package docker

import (
	"context"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type Controller struct {
	ctx    context.Context
	cli    *client.Client
	Logger *zap.Logger
}

func NewController(logger *zap.Logger) (*Controller, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	controller := &Controller{ctx: ctx, cli: cli, Logger: logger}

	return controller, nil
}
