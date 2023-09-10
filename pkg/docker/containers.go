package docker

import (
	"fmt"
	"reflect"
	"strings"

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

func (c *Controller) ContainerStop(containerId string) error {
	return c.cli.ContainerStop(c.ctx, containerId, container.StopOptions{})
}

func (c *Controller) ContainerInspect(containerId string) (types.ContainerJSON, error) {
	resp, err := c.cli.ContainerInspect(c.ctx, containerId)
	if err != nil {
		return types.ContainerJSON{}, err
	}

	return resp, nil
}

func (c *Controller) IsConfigDifferent(containerConfig, comparisonContainer types.ContainerJSON) bool {
	if comparisonContainer.Config.Image != containerConfig.Config.Image {
		return true
	}

	if comparisonContainer.HostConfig.NetworkMode != containerConfig.HostConfig.NetworkMode {
		return true
	}

	fmt.Println(comparisonContainer.HostConfig.RestartPolicy.Name, containerConfig.HostConfig.RestartPolicy.Name)
	if comparisonContainer.HostConfig.RestartPolicy.Name != containerConfig.HostConfig.RestartPolicy.Name {
		return true
	}

	if !reflect.DeepEqual(comparisonContainer.Config.Labels, containerConfig.Config.Labels) {
		return true
	}

	if !reflect.DeepEqual(comparisonContainer.HostConfig.Binds, containerConfig.HostConfig.Binds) {
		return true
	}

	if !reflect.DeepEqual(comparisonContainer.Config.Env, containerConfig.Config.Env) {
		return true
	}

	return false
}

func (c *Controller) ContainerSetupConfig(jobName string, config *pb.Container) *ContainerConfig {
	uuid := uuid.New()
	containerName := fmt.Sprintf("%s-%s", jobName, uuid.String())

	labels := make(map[string]string)

	labels[LABEL_DOKKUP_JOB_NAME] = jobName

	for _, label := range config.Labels {
		labelSplit := strings.Split(label, "=")
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

func (c *Controller) CreateTemporaryContainer(request *pb.DeployJobRequest) (string, error) {
	containerConfig := c.ContainerSetupConfig(request.Name, request.Container)

	container, err := c.ContainerCreate(request.Name+"-temporary", containerConfig.Config, containerConfig.HostConfig)
	if err != nil {
		return "", err
	}

	return container.ID, nil
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

func (c *Controller) AppendRollbackToContainers(containers []types.Container) error {
	for _, cont := range containers {
		err := c.cli.ContainerRename(c.ctx, cont.ID, cont.Names[0]+"-rollback")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) RemoveRollbackFromContainers(containers []types.Container) error {
	for _, cont := range containers {
		newName := strings.ReplaceAll(cont.Names[0], "-rollback", "")
		err := c.cli.ContainerRename(c.ctx, cont.ID, newName)
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

func (c *Controller) DeleteRollbackContainers(jobName string) error {
	rollbackContainers, err := c.GetRollbackContainers(jobName)
	if err != nil {
		return err
	}

	for _, container := range rollbackContainers {
		err = c.ContainerRemove(container.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) GetRollbackContainers(jobName string) ([]types.Container, error) {
	allContainers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	foundContainers := []types.Container{}
	for _, container := range allContainers {
		if container.Labels[LABEL_DOKKUP_JOB_NAME] == jobName && strings.Contains(container.Names[0], "-rollback") {
			foundContainers = append(foundContainers, container)
		}
	}

	return foundContainers, nil
}

func (c *Controller) ContainerRemove(containerId string) error {
	return c.cli.ContainerRemove(c.ctx, containerId, types.ContainerRemoveOptions{})
}

func (c *Controller) GetContainersByJobName(jobName string) ([]types.Container, error) {
	allContainers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	foundContainers := []types.Container{}
	for _, container := range allContainers {
		if container.Labels[LABEL_DOKKUP_JOB_NAME] == jobName && !strings.Contains(container.Names[0], "-temporary") && !strings.Contains(container.Names[0], "-rollback") {
			foundContainers = append(foundContainers, container)
		}
	}

	return foundContainers, nil
}
