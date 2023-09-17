package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/XiovV/dokkup/pkg/config"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (a *App) showNodeStatuses(inventory *config.Inventory, job *config.Job) error {
	group, ok := inventory.GetGroup(job.Group)
	if !ok {
		return fmt.Errorf("couldn't find group '%s'", job.Group)
	}

	fmt.Print("Node statuses:\n\n")
	nodeStatusesTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(nodeStatusesTable, "NAME\tSTATUS\tCONTAINERS\tUPDATE")

	var unavailableNodes int
	for _, nodeName := range group.Nodes {
		node, ok := inventory.GetNode(nodeName)
		if !ok {
			return fmt.Errorf("couldn't find node '%s", nodeName)
		}

		jobStatus, err := a.getNodeStatus(job.Name, node)
		if err != nil {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t", nodeName, err.Error(), 0, 0, false)
			fmt.Fprintln(nodeStatusesTable, out)

			unavailableNodes++
			continue
		}

		out := fmt.Sprintf("%s\t%s\t%d/%d\t%t", nodeName, NODE_STATUS_ONLINE, jobStatus.RunningContainers, jobStatus.TotalContainers, jobStatus.ShouldUpdate)
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

func (a *App) getNodeStatus(jobName string, node config.Node) (*pb.JobStatus, error) {
	jobStatus, err := a.pingNode(jobName, node)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			return nil, errors.New(NODE_STATUS_UNAUTHENTICATED)
		case codes.Unavailable:
			return nil, errors.New(NODE_STATUS_OFFLINE)
		default:
			return nil, err
		}
	}

	return jobStatus, nil
}

func (a *App) pingNode(jobName string, node config.Node) (*pb.JobStatus, error) {
	client, err := a.initClient(node.Location)
	if err != nil {
		return nil, fmt.Errorf("couldn't init connection: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", node.Key)

	jobStatus, err := client.GetJobStatus(ctx, &pb.GetJobStatusRequest{Name: jobName})
	if err != nil {
		return nil, fmt.Errorf("couldn't get node status: %w", err)
	}

	return jobStatus, nil
}

func (a *App) showStopJobSummaryTable(job *config.Job) {
	fmt.Print("Stop job summary:\n\n")

	jobSummaryTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(jobSummaryTable, "NAME\tIMAGE\tGROUP")

	out := fmt.Sprintf("%s\t%s\t%s\n", job.Name, job.Container[0].Image, job.Group)
	fmt.Fprintln(jobSummaryTable, out)

	jobSummaryTable.Flush()
}

func (a *App) showJobSummaryTable(job *config.Job) {
	fmt.Print("Deployment summary:\n\n")

	jobSummaryTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(jobSummaryTable, "NAME\tIMAGE\tRESTART\tCOUNT\tGROUP\tNETWORK")

	out := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\t%s\n", job.Name, job.Container[0].Image, job.Container[0].Restart, job.Count, job.Group, job.Container[0].Network)
	fmt.Fprintln(jobSummaryTable, out)

	jobSummaryTable.Flush()
}
