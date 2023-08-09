package docker

import (
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (c *Controller) ContainerCreate(containerName string, containerConfig *container.Config, hostConfig *container.HostConfig) error {
	resp, err := c.cli.ContainerCreate(c.ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	fmt.Println("container created successfully:", resp.ID)

	c.cli.ContainerStart(c.ctx, resp.ID, types.ContainerStartOptions{})
	return nil
}

func (c *Controller) ContainerInspect(containerId string) types.ContainerJSON {
	resp, err := c.cli.ContainerInspect(c.ctx, containerId)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
