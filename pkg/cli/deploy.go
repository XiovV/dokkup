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

	err = a.showNodeStatuses(inventory, job)
	if err != nil {
		log.Fatal("couldn't show node statuses: ", err)
	}

	shouldContinue, err := a.showConfirmationPrompt(ctx)
	if err != nil {
		log.Fatal("confirmation prompt error: ", err)
	}

	if !shouldContinue {
		return nil
	}

	fmt.Print("\n")

	err = a.deployJobs(job, inventory)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *App) deployJobs(job *config.Job, inventory *config.Inventory) error {
	group, ok := inventory.GetGroup(job.Group)
	if !ok {
		return fmt.Errorf("couldn't find group: %s", group.Name)
	}

	for _, nodeName := range group.Nodes {
		node, ok := inventory.GetNode(nodeName)
		if !ok {
			return fmt.Errorf("couldn't find node: %s", nodeName)
		}

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
		ports = append(ports, &pb.Port{In: port.In})
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
