package services

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

const PingCtxTimeout = 60 * time.Second

type ScsPingClient struct {
	grpcManager *grpcmanagers.PingGrpcManager
	grpcClient  pb.PingClient
}

func NewScsPingClient(request *models.PingClientRequest) (*ScsPingClient, momentoerrors.MomentoSvcErr) {
	pingManager, err := grpcmanagers.NewPingGrpcManager(&models.PingGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
	})
	if err != nil {
		return nil, err
	}
	return &ScsPingClient{
		grpcManager: pingManager,
		grpcClient:  pb.NewPingClient(pingManager.Conn),
	}, nil
}

func (client *ScsPingClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (client *ScsPingClient) Ping(ctx context.Context) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, PingCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.Ping(ctx, &pb.XPingRequest{})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}
