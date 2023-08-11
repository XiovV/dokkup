package docker

import (
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (c *Controller) ContainerCreate(containerName string, containerConfig *container.Config, hostConfig *container.HostConfig) (container.CreateResponse, error) {
	resp, err := c.cli.ContainerCreate(c.ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return container.CreateResponse{}, err
	}

	return resp, nil
}

func (c *Controller) ContainerStart(containerId string) error {
	return c.cli.ContainerStart(c.ctx, containerId, types.ContainerStartOptions{})
}

func (c *Controller) ContainerInspect(containerId string) types.ContainerJSON {
	resp, err := c.cli.ContainerInspect(c.ctx, containerId)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
