package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"github.com/XiovV/dokkup/pkg/config"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"github.com/gosuri/uilive"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	NODE_STATUS_ONLINE          = "ONLINE"
	NODE_STATUS_OFFLINE         = "OFFLINE"
	NODE_STATUS_UNAUTHENTICATED = "API KEY INVALID"
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

func (a *App) stopJobCmd(ctx *cli.Context) error {
	job, inventory, err := a.readJobAndInventory(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("This will stop a job called %s in the group %s", job.Name, job.Group)

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

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fmt.Println(resp)
	}

	return nil
}

func (a *App) readJobAndInventory(ctx *cli.Context) (*config.Job, *config.Inventory, error) {
	inventoryFlag := ctx.String("inventory")
	inventory, err := config.ReadInventory(inventoryFlag)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't read inventory: %w", err)
	}

	jobArg := ctx.Args().First()

	if len(jobArg) == 0 {
		return nil, nil, errors.New("please provide a job file")
	}

	job, err := config.ReadJob(jobArg)
	if err != nil {
		return nil, nil, err
	}

	return job, inventory, nil
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

func (a *App) jobToDeployJobRequest(job *config.Job) *pb.DeployJobRequest {
	count := int32(job.Count)
	container := job.Container[0]

	ports := []*pb.Port{}

	for _, port := range container.Ports {
		ports = append(ports, &pb.Port{In: port.In})
	}

	return &pb.DeployJobRequest{
		Count: count,
		Name:  job.Name,
		Container: &pb.Container{
			Image: container.Image,
			Ports: ports,
		},
	}
}

func (a *App) showNodeStatuses(inventory *config.Inventory, job *config.Job) error {
	group, ok := inventory.GetGroup(job.Group)
	if !ok {
		return fmt.Errorf("couldn't find group '%s'", job.Group)
	}

	var unavailableNodes int

	fmt.Print("Node statuses:\n\n")
	nodeStatusesTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(nodeStatusesTable, "NAME\tSTATUS")

	for _, nodeName := range group.Nodes {
		node, ok := inventory.GetNode(nodeName)
		if !ok {
			return fmt.Errorf("couldn't find node '%s", nodeName)
		}

		nodeStatus, err := a.getNodeStatus(node)
		if err != nil {
			log.Fatal("couldn't get node status: ", err)
		}

		if nodeStatus == NODE_STATUS_OFFLINE || nodeStatus == NODE_STATUS_UNAUTHENTICATED {
			unavailableNodes++
		}

		out := fmt.Sprintf("%s\t%s", nodeName, nodeStatus)
		fmt.Fprintln(nodeStatusesTable, out)
	}

	nodeStatusesTable.Flush()

	if unavailableNodes == 1 {
		fmt.Printf("\nWarning! It seems that there is %d unavailable node, it will be skipped.\n", unavailableNodes)
	}

	if unavailableNodes > 1 {
		fmt.Printf("\nWarning! It seems that there are %d unavailable nodes, they will be skipped.\n", unavailableNodes)
	}

	return nil
}

func (a *App) getNodeStatus(node config.Node) (string, error) {
	err := a.pingNode(node)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			return NODE_STATUS_UNAUTHENTICATED, nil
		case codes.Unavailable:
			return NODE_STATUS_OFFLINE, nil
		default:
			return "", err
		}
	}

	return NODE_STATUS_ONLINE, nil
}

func (a *App) pingNode(node config.Node) error {
	client, err := a.initClient(node.Location)
	if err != nil {
		return fmt.Errorf("couldn't init connection: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", node.Key)

	_, err = client.CheckAPIKey(ctx, &pb.CheckAPIKeyRequest{})
	if err != nil {
		return fmt.Errorf("couldn't check API key: %w", err)
	}

	return nil
}

func (a *App) showJobSummaryTable(job *config.Job) {
	fmt.Print("Deployment summary:\n\n")

	jobSummaryTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(jobSummaryTable, "NAME\tIMAGE\tRESTART\tCOUNT\tGROUP\tNETWORKS")

	out := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\t%s\n", job.Name, job.Container[0].Image, job.Container[0].Restart, job.Count, job.Group, job.Container[0].Networks)
	fmt.Fprintln(jobSummaryTable, out)

	jobSummaryTable.Flush()
}

func (a *App) showConfirmationPrompt(ctx *cli.Context) (bool, error) {
	if !ctx.Bool("yes") {
		fmt.Print("\nAre you sure you want to proceed? (y/n) ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		if input == "n\n" {
			return false, nil
		}
	}

	return true, nil
}
