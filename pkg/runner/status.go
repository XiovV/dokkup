package runner

import (
	"github.com/XiovV/dokkup/pkg/docker"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/docker/docker/api/types"
)

type JobContainers struct {
	TotalContainers         int
	RunningContainers       int
	RollbackContainers      int
	RunningContainerConfig  types.ContainerJSON
	RollbackContainerConfig types.ContainerJSON
}

func (j *JobRunner) GetJobStatus(job *pb.Job) (JobContainers, error) {
	totalContainers, err := j.Controller.GetContainers(job.Name, docker.GetContainersOptions{Stopped: true})
	if err != nil {
		return JobContainers{}, err
	}

	if len(totalContainers) == 0 {
		return JobContainers{
			TotalContainers:    0,
			RunningContainers:  0,
			RollbackContainers: 0,
		}, nil
	}

	containerConfig, err := j.Controller.ContainerInspect(totalContainers[0].ID)
	if err != nil {
		return JobContainers{}, err
	}

	runningContainers, err := j.Controller.GetContainers(job.Name, docker.GetContainersOptions{Stopped: false})
	if err != nil {
		return JobContainers{}, err
	}

	rollbackContainers, err := j.Controller.GetContainers(job.Name, docker.GetContainersOptions{Rollback: true})
	if err != nil {
		return JobContainers{}, err
	}

	if len(rollbackContainers) == 0 {
		return JobContainers{
			TotalContainers:        len(totalContainers),
			RunningContainers:      len(runningContainers),
			RunningContainerConfig: containerConfig,
		}, nil
	}

	rollbackContainerConfig, err := j.Controller.ContainerInspect(rollbackContainers[0].ID)
	if err != nil {
		return JobContainers{}, nil
	}

	return JobContainers{
		TotalContainers:         len(totalContainers),
		RunningContainers:       len(runningContainers),
		RollbackContainers:      len(rollbackContainers),
		RunningContainerConfig:  containerConfig,
		RollbackContainerConfig: rollbackContainerConfig,
	}, nil
}
