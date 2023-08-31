package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"go.uber.org/zap"

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

func (c *Controller) ContainerSetupConfig(jobName string, config *pb.Container) *ContainerConfig {
	uuid := uuid.New()
	containerName := fmt.Sprintf("%s-%s", jobName, uuid.String())

	labels := make(map[string]string)

	for _, l := range config.Labels {
		labelSplit := strings.Split(l, "=")
		labels[labelSplit[0]] = labelSplit[1]
	}

	containerConfig := &container.Config{
		Image:        config.Image,
		ExposedPorts: nat.PortSet{},
		Env:          config.Environment,
		Labels:       labels,
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{},
		Binds:        config.Volumes,
		RestartPolicy: container.RestartPolicy{
			Name: config.Restart,
		},
		NetworkMode: container.NetworkMode(config.Network),
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

func (c *Controller) CreateContainersFromRequest(request *pb.DeployJobRequest, stream pb.Dokkup_DeployJobServer) ([]string, error) {
	createdContainers := []string{}
	for i := 0; i < int(request.Count); i++ {
		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Configuring container (%d/%d)", i+1, request.Count)})

		containerConfig := c.ContainerSetupConfig(request.Name, request.Container)

		c.Logger.Info("attempting to create a container", zap.String("containerName", containerConfig.ContainerName))
		resp, err := c.ContainerCreate(containerConfig.ContainerName, containerConfig.Config, containerConfig.HostConfig)
		if err != nil {
			return nil, err
		}

		c.Logger.Info("container created successfully", zap.String("containerName", containerConfig.ContainerName))

		createdContainers = append(createdContainers, resp.ID)
	}

	return createdContainers, nil
}

func (c *Controller) StartContainers(containerIDs []string, stream pb.Dokkup_DeployJobServer) error {
	for i, container := range containerIDs {
		c.Logger.Info("attempting to start a container", zap.String("containerId", container))

		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Starting container (%d/%d)", i+1, len(containerIDs))})
		if err := c.ContainerStart(container); err != nil {
			return err
		}

		c.Logger.Info("container started successfully", zap.String("containerId", container))
	}

	return nil
}

func (c *Controller) ContainerDoesExist(containerName string) (bool, error) {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false, err
	}

	for _, container := range containers {
		if strings.Contains(container.Names[0], containerName) {
			return true, nil
		}
	}

	return false, nil
}

func (c *Controller) StopContainers(containers []types.Container, stream pb.Dokkup_StopJobServer) error {
	for i, cont := range containers {
		stream.Send(&pb.StopJobResponse{Message: fmt.Sprintf("Stopping container (%d/%d)", i+1, len(containers))})
		err := c.cli.ContainerStop(c.ctx, cont.ID, container.StopOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) DeleteContainers(containers []types.Container, stream pb.Dokkup_StopJobServer) error {
	for i, cont := range containers {
		stream.Send(&pb.StopJobResponse{Message: fmt.Sprintf("Deleting container (%d/%d)", i+1, len(containers))})
		err := c.cli.ContainerRemove(c.ctx, cont.ID, types.ContainerRemoveOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) GetContainersByJobName(jobName string) ([]types.Container, error) {
	allContainers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	foundContainers := []types.Container{}
	for _, container := range allContainers {
		if strings.Contains(container.Names[0], jobName) {
			foundContainers = append(foundContainers, container)
		}
	}

	return foundContainers, nil
}
