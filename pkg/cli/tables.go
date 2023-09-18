package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/XiovV/dokkup/pkg/config"
)

type JobStatus struct {
	Node              config.Node
	NodeStatus        string
	RunningContainers int
	TotalContainers   int
	ShouldUpdate      bool
}

func (a *App) showDeployJobStatuses(jobStatuses []JobStatus) error {
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

		nodeName := jobStatus.Node.Name

		if jobStatus.TotalContainers == 0 {
			nodeName += "*"
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
