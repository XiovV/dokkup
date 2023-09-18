package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/XiovV/dokkup/pkg/config"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JobStatus struct {
	Node              config.Node
	NodeStatus        string
	RunningContainers int
	TotalContainers   int
	ShouldUpdate      bool
}

func (a *App) showJobStatuses(jobStatuses []JobStatus) error {
	fmt.Print("Node statuses:\n\n")
	nodeStatusesTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(nodeStatusesTable, "NAME\tSTATUS\tCONTAINERS\tUPDATE")

	var unavailableNodes int
	for _, jobStatus := range jobStatuses {
		if jobStatus.NodeStatus == NODE_STATUS_OFFLINE || jobStatus.NodeStatus == NODE_STATUS_UNAUTHENTICATED {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t", jobStatus.Node.Name, jobStatus.NodeStatus, 0, 0, false)
			fmt.Fprintln(nodeStatusesTable, out)

			unavailableNodes++
			continue
		}

		out := fmt.Sprintf("%s\t%s\t%d/%d\t%t", jobStatus.Node.Name, NODE_STATUS_ONLINE, jobStatus.RunningContainers, jobStatus.TotalContainers, jobStatus.ShouldUpdate)
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

func (a *App) getJobStatus(job *config.Job, node config.Node) (JobStatus, error) {
	jobStatusResponse, err := a.pingNode(job, node)

	jobStatus := JobStatus{
		Node:       node,
		NodeStatus: NODE_STATUS_ONLINE,
	}

	if err != nil {
		jobStatus.RunningContainers = 0
		jobStatus.TotalContainers = 0
		jobStatus.ShouldUpdate = false

		switch status.Code(err) {
		case codes.Unauthenticated:
			jobStatus.NodeStatus = NODE_STATUS_UNAUTHENTICATED
			return jobStatus, nil
		case codes.Unavailable:
			jobStatus.NodeStatus = NODE_STATUS_OFFLINE
			return jobStatus, nil
		default:
			return JobStatus{}, err
		}
	}

	jobStatus.RunningContainers = int(jobStatusResponse.RunningContainers)
	jobStatus.TotalContainers = int(jobStatusResponse.TotalContainers)
	jobStatus.ShouldUpdate = jobStatusResponse.ShouldUpdate

	return jobStatus, nil
}

func (a *App) pingNode(job *config.Job, node config.Node) (*pb.JobStatus, error) {
	client, err := a.initClient(node.Location)
	if err != nil {
		return nil, fmt.Errorf("couldn't init connection: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", node.Key)

	request := a.jobToDeployJobRequest(job)

	jobStatus, err := client.GetJobStatus(ctx, request)
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
