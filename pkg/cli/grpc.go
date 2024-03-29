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
	if err != nil {
		jobStatus := JobStatus{Node: node}
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

	jobStatus := JobStatus{
		Node:              node,
		NodeStatus:        NODE_STATUS_ONLINE,
		RunningContainers: int(jobStatusResponse.RunningContainers),
		TotalContainers:   int(jobStatusResponse.TotalContainers),
		ShouldUpdate:      jobStatusResponse.ShouldUpdate,
		CanRollback:       jobStatusResponse.CanRollback,
		CurrentVersion:    jobStatusResponse.CurrentVersion,
		NewVersion:        jobStatusResponse.NewVersion,
		OldVersion:        jobStatusResponse.OldVersion,
		Containers:        jobStatusResponse.Containers,
		Image:             jobStatusResponse.Image,
	}

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
