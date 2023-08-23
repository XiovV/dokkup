package docker

import (
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	pb "github.com/XiovV/dokkup/pkg/grpc"
)

type ContainerConfig struct {
	ContainerName string
	Config        *container.Config
	HostConfig    *container.HostConfig
}

func (c *Controller) IsConfigDifferent(request *pb.DeployJobRequest) (bool, error) {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false, nil
	}

	var container types.Container
	for _, c := range containers {
		if strings.Contains(c.Names[0], request.Name) {
			container = c
		}

		return false, errors.New("couldn't find container")
	}

	config := c.ContainerInspect(container.ID)

	if request.Container.Image != config.Image {
		return true, nil
	}

	return false, nil
}
