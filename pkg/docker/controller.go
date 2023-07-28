package docker

import (
	"context"

	"github.com/docker/docker/client"
)

type Controller struct {
  ctx context.Context
  cli *client.Client
}

func NewController() (*Controller, error) {
  ctx := context.Background()
  cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    return nil, err
  }

  controller := &Controller{ctx: ctx, cli: cli}

  return controller, nil
}
