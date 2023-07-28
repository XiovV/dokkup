package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
)

func (c *Controller) ContainerCreate(containerName, containerImage string) error {
  resp, err := c.cli.ContainerCreate(c.ctx, &container.Config{Image: containerImage}, nil, nil, nil, containerName)
  if err != nil {
    return err
  }

  fmt.Println("container created successfully:", resp.ID)
  return nil
}
