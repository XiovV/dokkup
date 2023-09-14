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

func (a *App) stopJobCmd(ctx *cli.Context) error {
	job, inventory, err := a.readJobAndInventory(ctx)
	if err != nil {
		log.Fatal(err)
	}

	a.showStopJobSummaryTable(job)

	a.showNodeStatuses(inventory, job)

	shouldContinue, err := a.showConfirmationPrompt(ctx)
	if err != nil {
		log.Fatal("confirmation prompt error: ", err)
	}

	if !shouldContinue {
		return nil
	}

	fmt.Print("\n")

	err = a.stopJobs(job, inventory)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *App) stopJobs(job *config.Job, inventory *config.Inventory) error {
	group, ok := inventory.GetGroup(job.Group)
	if !ok {
		return fmt.Errorf("couldn't find group: %s", group.Name)
	}

	for _, nodeName := range group.Nodes {
		node, ok := inventory.GetNode(nodeName)
		if !ok {
			return fmt.Errorf("couldn't find node: %s", nodeName)
		}

		err := a.stopJob(node, job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) stopJob(node config.Node, job *config.Job) error {
	client, err := a.initClient(node.Location)
	if err != nil {
		return err
	}

	ctx, cancel := a.newAuthorizationContext(node.Key)
	defer cancel()

	request := &pb.StopJobRequest{
		Name: job.Name,
	}

	stream, err := client.StopJob(ctx, request)
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
