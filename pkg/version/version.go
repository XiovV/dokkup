package version

import (
	"fmt"

	"github.com/cnf/structhash"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type containerConfig struct {
	Image         string
	RestartPolicy string
	NetworkMode   container.NetworkMode
	Env           []string
	Labels        map[string]string
	Binds         []string
	ExposedPorts  nat.PortSet
}

func Hash(container types.ContainerJSON) string {
	cont := containerConfig{
		Image:         container.Image,
		RestartPolicy: container.HostConfig.RestartPolicy.Name,
		NetworkMode:   container.HostConfig.NetworkMode,
		Env:           container.Config.Env,
		Labels:        container.Config.Labels,
		Binds:         container.HostConfig.Binds,
		ExposedPorts:  container.Config.ExposedPorts,
	}

	versionHash := structhash.Md5(cont, 0)
	return fmt.Sprintf("%x", versionHash)
}
