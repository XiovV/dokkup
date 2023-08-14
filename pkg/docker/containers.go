package docker

import (
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"

	pb "github.com/XiovV/dokkup/pkg/grpc"
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

func (c *Controller) ContainerSetupConfig(config *pb.Container) *ContainerConfig {
	uuid := uuid.New()
	containerName := fmt.Sprintf("%s-%s", config.Name, uuid.String())

	containerConfig := &container.Config{
		Image:        config.Image,
		ExposedPorts: nat.PortSet{},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{},
	}

	for _, port := range config.Ports {
		internalPort := fmt.Sprintf("%s/tcp", port.In)
		containerConfig.ExposedPorts[nat.Port(internalPort)] = struct{}{}

		hostConfig.PortBindings[nat.Port(internalPort)] = []nat.PortBinding{{HostIP: "0.0.0.0"}}
	}

	return &ContainerConfig{
		ContainerName: containerName,
		Config:        containerConfig,
		HostConfig:    hostConfig,
	}
}

func (c *Controller) CreateContainersFromRequest(request *pb.DeployJobRequest) ([]string, error) {
	createdContainers := []string{}
	for i := 0; i < int(request.Count); i++ {
		containerConfig := c.ContainerSetupConfig(request.Container)

		resp, err := c.ContainerCreate(containerConfig.ContainerName, containerConfig.Config, containerConfig.HostConfig)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		createdContainers = append(createdContainers, resp.ID)
	}

	return createdContainers, nil
}

func (c *Controller) StartContainers(containerIDs []string) error {
	for _, container := range containerIDs {
		if err := c.ContainerStart(container); err != nil {
			return err
		}
	}

	return nil
}
