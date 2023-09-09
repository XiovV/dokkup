package runner

import (
	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/docker/docker/api/types"
)

type JobRunner struct {
	Controller *docker.Controller
}

func NewJobRunner(controller *docker.Controller) *JobRunner {
	return &JobRunner{Controller: controller}
}

func (j *JobRunner) ShouldUpdateJob(jobName string, comparisonContainer types.ContainerJSON) (bool, error) {
	runningContainers, err := j.Controller.GetContainersByJobName(jobName)
	if err != nil {
		return false, err
	}

	if len(runningContainers) == 0 {
		return false, nil
	}

	containerConfig, err := j.Controller.ContainerInspect(runningContainers[0].ID)
	if err != nil {
		return false, err
	}

	isDifferent := j.Controller.IsConfigDifferent(containerConfig, comparisonContainer)
	if !isDifferent {
		return false, nil
	}

	return true, nil
}

func (j *JobRunner) DoesJobExist(jobName string) bool {
	containers, err := j.Controller.GetContainersByJobName(jobName)
	if err != nil {
		return false
	}

	return len(containers) != 0
}

func (j *JobRunner) RunDeployment(stream pb.Dokkup_DeployJobServer, request *pb.DeployJobRequest) error {
	createdContainers, err := j.Controller.CreateContainersFromRequest(request, stream)
	if err != nil {
		return err
	}

	return j.Controller.StartContainers(createdContainers, stream)
}

func (j *JobRunner) RunUpdate(request *pb.DeployJobRequest, stream pb.Dokkup_DeployJobServer) error {
	err := j.Controller.DeleteRollbackContainers(request.Name)
	if err != nil {
		return err
	}

	oldContainers, err := j.Controller.GetContainersByJobName(request.Name)
	if err != nil {
		return err
	}

	err = j.Controller.AppendRollbackToContainers(oldContainers)
	if err != nil {
		return err
	}

	newContainers, err := j.Controller.CreateContainersFromRequest(request, stream)
	if err != nil {
		return err
	}

	err = j.Controller.StopContainers(oldContainers)
	if err != nil {
		return err
	}

	err = j.Controller.StartContainers(newContainers, stream)
	if err != nil {
		j.abortUpdate(request.Name)
		return err
	}

	return nil
}

func (j *JobRunner) abortUpdate(jobName string) error {
	rollbackContainers, err := j.Controller.GetRollbackContainers(jobName)
	if err != nil {
		return err
	}

	for _, container := range rollbackContainers {
		err := j.Controller.ContainerStart(container.ID)
		if err != nil {
			return err
		}
	}

	oldContainers, err := j.Controller.GetContainersByJobName(jobName)
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
