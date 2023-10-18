package runner

import (
	"fmt"

	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/docker/docker/api/types"
	"go.uber.org/zap"
)

type JobRunner struct {
	Controller *docker.Controller
	Logger     *zap.Logger
}

func NewJobRunner(controller *docker.Controller, logger *zap.Logger) *JobRunner {
	return &JobRunner{Controller: controller, Logger: logger}
}

func (j *JobRunner) StartContainers(containers []types.Container, stream pb.Dokkup_DeployJobServer) error {
	for i, container := range containers {
		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Starting container (%d/%d)", i+1, len(containers))})
		err := j.Controller.ContainerStart(container.ID)
		if err != nil {
			j.Logger.Error("could not start container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	return nil
}

func (j *JobRunner) ShouldUpdateJob(request *pb.Job, currentContainers []types.Container) (bool, error) {
	if len(currentContainers) == 0 {
		return false, nil
	}

	j.Logger.Debug("creating temporary container")
	temporaryContainer, err := j.Controller.CreateTemporaryContainer(request)
	if err != nil {
		j.Logger.Error("failed to create temporary container", zap.Error(err))
		return false, err
	}

	containerConfig, err := j.Controller.ContainerInspect(currentContainers[0].ID)
	if err != nil {
		return false, err
	}

	isDifferent := j.Controller.IsConfigDifferent(containerConfig, temporaryContainer)
	if !isDifferent {
		return false, nil
	}

	return true, nil
}

func (j *JobRunner) DoesJobExist(jobName string) bool {
	containers, err := j.Controller.GetContainers(jobName, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		return false
	}

	return len(containers) != 0
}

func (j *JobRunner) UpscaleJob(count int, request *pb.Job, stream pb.Dokkup_DeployJobServer) error {
	createdContainers, err := j.createContainersFromRequest(count, request, stream)
	if err != nil {
		return err
	}

	return j.startContainers(createdContainers, stream)
}

func (j *JobRunner) DownscaleJob(count int, currentContainers []types.Container, stream pb.Dokkup_DeployJobServer) error {
	j.Logger.Info("stopping job containers")
	for i := 0; i < count; i++ {
		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Stopping container (%d/%d)", i+1, count)})
		err := j.Controller.ContainerStop(currentContainers[i].ID)
		if err != nil {
			j.Logger.Error("could not stop container", zap.Error(err), zap.String("containerId", currentContainers[i].ID))
			return err
		}
	}

	j.Logger.Info("deleting job containers")
	for i := 0; i < count; i++ {
		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Removing container (%d/%d)", i+1, count)})
		err := j.Controller.ContainerRemove(currentContainers[i].ID)
		if err != nil {
			j.Logger.Error("could not remove job container", zap.Error(err), zap.String("containerId", currentContainers[i].ID))
			return err
		}
	}

	return nil
}

func (j *JobRunner) RunDeployment(stream pb.Dokkup_DeployJobServer, request *pb.Job) error {
	j.Logger.Debug("attempting to pull image", zap.String("image", request.Container.Image))

	stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to pull image: %s", request.Container.Image)})
	err := j.Controller.ImagePull(request.Container.Image)
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	createdContainers, err := j.createContainersFromRequest(int(request.Count), request, stream)
	if err != nil {
		return err
	}

	return j.startContainers(createdContainers, stream)
}

func (j *JobRunner) createContainersFromRequest(count int, request *pb.Job, stream pb.Dokkup_DeployJobServer) ([]string, error) {
	createdContainers := []string{}
	for i := 0; i < count; i++ {
		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Configuring container (%d/%d)", i+1, request.Count)})
		containerConfig := j.Controller.ContainerSetupConfig(request.Name, request.Container)

		resp, err := j.Controller.ContainerCreate(containerConfig.ContainerName, containerConfig.Config, containerConfig.HostConfig)
		if err != nil {
			return nil, err
		}

		createdContainers = append(createdContainers, resp.ID)
	}

	return createdContainers, nil
}

func (j *JobRunner) StopJob(request *pb.StopJobRequest, stream pb.Dokkup_StopJobServer) error {
	j.Logger.Debug("getting job containers")
	containers, err := j.Controller.GetContainers(request.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		j.Logger.Error("could not get containers", zap.Error(err))
		return err
	}

	if len(containers) == 0 {
		stream.Send(&pb.StopJobResponse{Message: "No running containers available"})
		return nil
	}

	j.Logger.Info("stopping job containers")
	for i, container := range containers {
		stream.Send(&pb.StopJobResponse{Message: fmt.Sprintf("Stopping container (%d/%d)", i+1, len(containers))})
		err := j.Controller.ContainerStop(container.ID)
		if err != nil {
			j.Logger.Error("could not stop container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	if !request.Purge {
		return nil
	}

	j.Logger.Info("deleting job containers")
	for i, container := range containers {
		stream.Send(&pb.StopJobResponse{Message: fmt.Sprintf("Removing container (%d/%d)", i+1, len(containers))})
		err := j.Controller.ContainerRemove(container.ID)
		if err != nil {
			j.Logger.Error("could not remove job container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	j.Logger.Debug("getting rollback containers")
	rollbackContainers, err := j.Controller.GetContainers(request.Name, docker.GetContainersOptions{Rollback: true})
	if err != nil {
		j.Logger.Error("could not get rollback containers", zap.Error(err))
		return err
	}

	if len(rollbackContainers) == 0 {
		return nil
	}

	j.Logger.Info("deleting rollback containers")
	for i, container := range rollbackContainers {
		stream.Send(&pb.StopJobResponse{Message: fmt.Sprintf("Removing rollback container (%d/%d)", i+1, len(containers))})
		err := j.Controller.ContainerRemove(container.ID)
		if err != nil {
			j.Logger.Error("could not remove rollback container", zap.Error(err), zap.String("containerId", container.ID))
			return err
		}
	}

	stream.Send(&pb.StopJobResponse{Message: "Job stopped successfully"})

	return nil
}

func (j *JobRunner) RunUpdate(request *pb.Job, stream pb.Dokkup_DeployJobServer) error {
	j.Logger.Debug("attempting to pull image", zap.String("image", request.Container.Image))

	stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Attempting to pull image: %s", request.Container.Image)})
	err := j.Controller.ImagePull(request.Container.Image)
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	j.Logger.Debug("removing previous rollback containers")
	err = j.Controller.DeleteRollbackContainers(request.Name)
	if err != nil {
		j.Logger.Error("could not remove previous rollback containers", zap.Error(err))
		return err
	}

	oldContainers, err := j.Controller.GetContainers(request.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		return err
	}

	j.Logger.Debug("setting rollback containers")
	err = j.Controller.AppendRollbackToContainers(oldContainers)
	if err != nil {
		j.Logger.Error("failed to set rollback containers", zap.Error(err))
		return err
	}

	j.Logger.Debug("creating new containers")
	newContainers, err := j.createContainersFromRequest(int(request.Count), request, stream)
	if err != nil {
		j.Logger.Error("failed to create new containers", zap.Error(err))
		return err
	}

	err = j.stopContainers(oldContainers, stream)
	if err != nil {
		return err
	}

	err = j.startContainers(newContainers, stream)
	if err != nil {
		j.abortUpdate(request.Name)
		return err
	}

	stream.Send(&pb.DeployJobResponse{Message: "Update successful"})
	return nil
}

func (j *JobRunner) startContainers(containerIDs []string, stream pb.Dokkup_DeployJobServer) error {
	for i, container := range containerIDs {
		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Starting container (%d/%d)", i+1, len(containerIDs))})
		err := j.Controller.ContainerStart(container)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JobRunner) stopContainers(containers []types.Container, stream pb.Dokkup_DeployJobServer) error {
	for i, container := range containers {
		stream.Send(&pb.DeployJobResponse{Message: fmt.Sprintf("Stopping container (%d/%d)", i+1, len(containers))})
		err := j.Controller.ContainerStop(container.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JobRunner) abortUpdate(jobName string) error {
	rollbackContainers, err := j.Controller.GetContainers(jobName, docker.GetContainersOptions{Rollback: true})
	if err != nil {
		return err
	}

	for _, container := range rollbackContainers {
		err := j.Controller.ContainerStart(container.ID)
		if err != nil {
			return err
		}
	}

	oldContainers, err := j.Controller.GetContainers(jobName, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		return err
	}

	for _, container := range oldContainers {
		err := j.Controller.ContainerRemove(container.ID)
		if err != nil {
			return err
		}
	}

	err = j.Controller.RemoveRollbackFromContainers(rollbackContainers)
	if err != nil {
		return nil
	}

	return nil
}
