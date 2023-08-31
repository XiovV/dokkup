package cli

import (
	"context"
	"fmt"
	"log"
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
