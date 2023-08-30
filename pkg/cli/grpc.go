package cli

import (
	"context"

	pb "github.com/XiovV/dokkup/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
