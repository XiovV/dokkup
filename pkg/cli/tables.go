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
	CanRollback       bool
	CurrentVersion    string
	NewVersion        string
	OldVersion        string
}

func (a *App) showDeployJobStatuses(jobStatuses []JobStatus, job *config.Job) error {
	fmt.Print("Node statuses:\n\n")
	nodeStatusesTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(nodeStatusesTable, "NAME\tSTATUS\tCONTAINERS\tUPDATE\tVERSION")

	var unavailableNodes int
	for _, jobStatus := range jobStatuses {
		if jobStatus.NodeStatus == NODE_STATUS_OFFLINE || jobStatus.NodeStatus == NODE_STATUS_UNAUTHENTICATED {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t\t%s", jobStatus.Node.Name, jobStatus.NodeStatus, 0, 0, false, "")
			fmt.Fprintln(nodeStatusesTable, out)

			unavailableNodes++
			continue
		}

		nodeName := jobStatus.Node.Name

		if jobStatus.TotalContainers == 0 && jobStatus.ShouldUpdate {
			nodeName += "*"

			out := fmt.Sprintf("%s\t%s\t%d -> %d\t%t\t%s", nodeName, NODE_STATUS_ONLINE, 0, job.Count, jobStatus.ShouldUpdate, jobStatus.NewVersion[:7])
			fmt.Fprintln(nodeStatusesTable, out)
			continue
		}

		if jobStatus.ShouldUpdate {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t\t%s -> %s", nodeName, NODE_STATUS_ONLINE, jobStatus.RunningContainers, jobStatus.TotalContainers, jobStatus.ShouldUpdate, jobStatus.CurrentVersion[:7], jobStatus.NewVersion[:7])
			fmt.Fprintln(nodeStatusesTable, out)
			continue
		}

		out := fmt.Sprintf("%s\t%s\t%d/%d\t%t\t%s", nodeName, NODE_STATUS_ONLINE, jobStatus.RunningContainers, jobStatus.TotalContainers, jobStatus.ShouldUpdate, jobStatus.CurrentVersion[:7])
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

func (a *App) showJobSummaryTable(job *config.Job) {
	fmt.Print("Deployment summary:\n\n")

	jobSummaryTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(jobSummaryTable, "NAME\tIMAGE\tRESTART\tCOUNT\tGROUP\tNETWORK")

	out := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\t%s\n", job.Name, job.Container[0].Image, job.Container[0].Restart, job.Count, job.Group, job.Container[0].Network)
	fmt.Fprintln(jobSummaryTable, out)

	jobSummaryTable.Flush()
}

func (a *App) showStopJobSummaryTable(job *config.Job) {
	fmt.Print("Stop job summary:\n\n")

	jobSummaryTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(jobSummaryTable, "NAME\tIMAGE\tGROUP")

	out := fmt.Sprintf("%s\t%s\t%s\n", job.Name, job.Container[0].Image, job.Group)
	fmt.Fprintln(jobSummaryTable, out)

	jobSummaryTable.Flush()
}

func (a *App) showStopJobStatuses(jobStatuses []JobStatus, shouldPurge bool) error {
	fmt.Print("Node statuses:\n\n")
	nodeStatusesTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(nodeStatusesTable, "NAME\tSTATUS\tCONTAINERS\tPURGE")

	var unavailableNodes int
	for _, jobStatus := range jobStatuses {
		if jobStatus.NodeStatus == NODE_STATUS_OFFLINE || jobStatus.NodeStatus == NODE_STATUS_UNAUTHENTICATED {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t", jobStatus.Node.Name, jobStatus.NodeStatus, 0, 0, false)
			fmt.Fprintln(nodeStatusesTable, out)

			unavailableNodes++
			continue
		}

		out := fmt.Sprintf("%s\t%s\t%d -> %d\t%t", jobStatus.Node.Name, NODE_STATUS_ONLINE, jobStatus.RunningContainers, 0, shouldPurge)
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

func (a *App) showRollbackJobStatuses(jobStatuses []JobStatus) error {
	fmt.Print("Node statuses:\n\n")
	nodeStatusesTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(nodeStatusesTable, "NAME\tSTATUS\tCONTAINERS\tROLLBACK\tVERSION")

	var unavailableNodes int
	for _, jobStatus := range jobStatuses {
		if jobStatus.NodeStatus == NODE_STATUS_OFFLINE || jobStatus.NodeStatus == NODE_STATUS_UNAUTHENTICATED {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t\t%s", jobStatus.Node.Name, jobStatus.NodeStatus, 0, 0, false, "")
			fmt.Fprintln(nodeStatusesTable, out)

			unavailableNodes++
			continue
		}

		nodeName := jobStatus.Node.Name

		if jobStatus.CanRollback {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t\t%s -> %s", nodeName, NODE_STATUS_ONLINE, jobStatus.RunningContainers, jobStatus.TotalContainers, jobStatus.CanRollback, jobStatus.CurrentVersion[:7], jobStatus.OldVersion[:7])
			fmt.Fprintln(nodeStatusesTable, out)

			continue
		}

		if jobStatus.TotalContainers == 0 {
			out := fmt.Sprintf("%s\t%s\t%d/%d\t%t\t%s", nodeName, NODE_STATUS_ONLINE, jobStatus.RunningContainers, jobStatus.TotalContainers, jobStatus.CanRollback, "")
			fmt.Fprintln(nodeStatusesTable, out)

			continue
		}

		out := fmt.Sprintf("%s\t%s\t%d/%d\t%t\t%s", nodeName, NODE_STATUS_ONLINE, jobStatus.RunningContainers, jobStatus.TotalContainers, jobStatus.CanRollback, jobStatus.CurrentVersion[:7])
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
