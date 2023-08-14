package docker

import "github.com/docker/docker/api/types/container"

type ContainerConfig struct {
	ContainerName string
	Config        *container.Config
	HostConfig    *container.HostConfig
}
