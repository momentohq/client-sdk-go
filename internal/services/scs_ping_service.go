package services

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ScsPingClient struct {
	requestTimeout time.Duration
	grpcManager    *grpcmanagers.PingGrpcManager
	grpcClient     pb.PingClient
}

func NewScsPingClient(request *models.PingClientRequest) (*ScsPingClient, momentoerrors.MomentoSvcErr) {
	pingManager, err := grpcmanagers.NewPingGrpcManager(&models.PingGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  request.Configuration.GetTransportStrategy().GetGrpcConfig(),
	})
	if err != nil {
		return nil, err
	}
	return &ScsPingClient{
		requestTimeout: request.Configuration.GetClientSideTimeout(),
		grpcManager:    pingManager,
		grpcClient:     pb.NewPingClient(pingManager.Conn),
	}, nil
}

func (client *ScsPingClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (client *ScsPingClient) Ping(ctx context.Context) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	requestMetadata := internal.CreateMetadata(ctx, internal.Ping)
	defer cancel()
	_, err := client.grpcClient.Ping(requestMetadata, &pb.XPingRequest{})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}
