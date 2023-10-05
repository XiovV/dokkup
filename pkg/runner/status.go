package runner

import (
	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/docker/docker/api/types"
)

type JobStatus struct {
	TotalContainers         []types.Container
	RunningContainers       []types.Container
	RollbackContainers      []types.Container
	RunningContainerConfig  types.ContainerJSON
	RollbackContainerConfig types.ContainerJSON
	Image                   string
}

func (j *JobRunner) GetJobStatus(job *pb.Job) (JobStatus, error) {
	totalContainers, err := j.Controller.GetContainers(job.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		return JobStatus{}, err
	}

	if len(totalContainers) == 0 {
		return JobStatus{}, nil
	}

	containerConfig, err := j.Controller.ContainerInspect(totalContainers[0].ID)
	if err != nil {
		return JobStatus{}, err
	}

	runningContainers, err := j.Controller.GetContainers(job.Name, docker.GetContainersOptions{Stopped: false})
	if err != nil {
		return JobStatus{}, err
	}

	rollbackContainers, err := j.Controller.GetContainers(job.Name, docker.GetContainersOptions{Rollback: true})
	if err != nil {
		return JobStatus{}, err
	}

	if len(rollbackContainers) == 0 {
		return JobStatus{
			TotalContainers:        totalContainers,
			RunningContainers:      runningContainers,
			RunningContainerConfig: containerConfig,
			Image:                  containerConfig.Config.Image,
		}, nil
	}

	rollbackContainerConfig, err := j.Controller.ContainerInspect(rollbackContainers[0].ID)
	if err != nil {
		return JobStatus{}, nil
	}

	return JobStatus{
		TotalContainers:         totalContainers,
		RunningContainers:       runningContainers,
		RollbackContainers:      rollbackContainers,
		RunningContainerConfig:  containerConfig,
		RollbackContainerConfig: rollbackContainerConfig,
		Image:                   containerConfig.Image,
	}, nil
}
