package cli

import (
	"fmt"
	"io"
	"log"

	"github.com/XiovV/dokkup/pkg/config"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/gosuri/uilive"
	"github.com/urfave/cli/v2"
)

func (a *App) jobCmd(ctx *cli.Context) error {
	job, inventory, err := a.readJobAndInventory(ctx)
	if err != nil {
		log.Fatal(err)
	}

	a.showJobSummaryTable(job)

	jobStatuses, err := a.getJobStatuses(job, inventory)
	if err != nil {
		log.Fatal(err)
	}

	a.showDeployJobStatuses(jobStatuses, job)

	shouldContinue, err := a.showConfirmationPrompt(ctx)
	if err != nil {
		log.Fatal("confirmation prompt error: ", err)
	}

	if !shouldContinue {
		return nil
	}

	fmt.Print("\n")

	availableNodes := a.extractAvailableNodes(jobStatuses)

	err = a.deployJobs(availableNodes, job)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *App) extractAvailableNodes(jobStatuses []JobStatus) []config.Node {
	availableNodes := []config.Node{}
	for _, status := range jobStatuses {
		if status.NodeStatus == NODE_STATUS_OFFLINE || status.NodeStatus == NODE_STATUS_UNAUTHENTICATED {
			continue
		}

		availableNodes = append(availableNodes, status.Node)
	}

	return availableNodes
}

func (a *App) getJobStatuses(job *config.Job, inventory *config.Inventory) ([]JobStatus, error) {
	group, ok := inventory.GetGroup(job.Group)
	if !ok {
		return nil, fmt.Errorf("couldn't find group '%s'", job.Group)
	}

	jobStatuses := []JobStatus{}
	for _, nodeName := range group.Nodes {
		node, ok := inventory.GetNode(nodeName)
		if !ok {
			return nil, fmt.Errorf("couldn't find node '%s", nodeName)
		}

		jobStatus, err := a.getJobStatus(job, node)
		if err != nil {
			return nil, err
		}

		jobStatuses = append(jobStatuses, jobStatus)
	}

	return jobStatuses, nil
}

func (a *App) deployJobs(nodes []config.Node, job *config.Job) error {
	for _, node := range nodes {
		err := a.deployJob(node, job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) deployJob(node config.Node, job *config.Job) error {
	client, err := a.initClient(node.Location)
	if err != nil {
		return err
	}

	ctx, cancel := a.newAuthorizationContext(node.Key)
	defer cancel()

	request := a.jobToDeployJobRequest(job)

	stream, err := client.DeployJob(ctx, request)
	if err != nil {
		return err
	}

	writer := uilive.New()
	writer.Start()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			writer.Stop()
			break
		}

		if err != nil {
			return err
		}

		fmt.Fprintf(writer, "%s: %s\n", node.Name, resp.GetMessage())
	}

	return nil
}

func (a *App) jobToDeployJobRequest(job *config.Job) *pb.Job {
	count := int32(job.Count)
	container := job.Container[0]

	ports := []*pb.Port{}

	for _, port := range container.Ports {
		ports = append(ports, &pb.Port{In: port.In, Out: port.Out})
	}

	return &pb.Job{
		Count: count,
		Name:  job.Name,
		Container: &pb.Container{
			Image:       container.Image,
			Ports:       ports,
			Network:     container.Network,
			Volumes:     container.Volumes,
			Environment: container.Environment,
			Restart:     container.Restart,
			Labels:      container.Labels,
		},
	}
}
