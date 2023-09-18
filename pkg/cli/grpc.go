package cli

import (
	"context"
	"fmt"

	"github.com/XiovV/dokkup/pkg/config"
	pb "github.com/XiovV/dokkup/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (a *App) initClient(target string) (pb.DokkupClient, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewDokkupClient(conn)
	return client, nil
}

func (a *App) newAuthorizationContext(nodeKey string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", nodeKey)

	return ctx, cancel
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
