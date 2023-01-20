package services

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/internal/utility"
)

const ControlPort = ":443"
const ControlCtxTimeout = 60 * time.Second

type ScsControlClient struct {
	grpcManager *grpcmanagers.ScsControlGrpcManager
	grpcClient  pb.ScsControlClient
}

func NewScsControlClient(request *models.ControlClientRequest) (*ScsControlClient, momentoerrors.MomentoSvcErr) {
	controlManager, err := grpcmanagers.NewScsControlGrpcManager(&models.ControlGrpcManagerRequest{
		AuthToken: request.AuthToken,
		Endpoint:  fmt.Sprint(request.Endpoint, ControlPort),
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &ScsControlClient{grpcManager: controlManager, grpcClient: pb.NewScsControlClient(controlManager.Conn)}, nil
}

func (client *ScsControlClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (client *ScsControlClient) CreateCache(ctx context.Context, request *models.CreateCacheRequest) momentoerrors.MomentoSvcErr {
	if !utility.IsCacheNameValid(request.CacheName) {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.CreateCache(ctx, &pb.XCreateCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsControlClient) CreateTopic(ctx context.Context, request *models.CreateTopicRequest) momentoerrors.MomentoSvcErr {
	if !utility.IsCacheNameValid(request.TopicName) {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Topic name cannot be empty", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.CreateCache(ctx, &pb.XCreateCacheRequest{CacheName: "topic-" + request.TopicName})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsControlClient) DeleteTopic(ctx context.Context, request *models.DeleteTopicRequest) momentoerrors.MomentoSvcErr {
	if !utility.IsCacheNameValid(request.TopicName) {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Topic name cannot be empty", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.DeleteCache(ctx, &pb.XDeleteCacheRequest{CacheName: "topic-" + request.TopicName})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsControlClient) DeleteCache(ctx context.Context, request *models.DeleteCacheRequest) momentoerrors.MomentoSvcErr {
	if !utility.IsCacheNameValid(request.CacheName) {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.DeleteCache(ctx, &pb.XDeleteCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsControlClient) ListCaches(ctx context.Context, request *models.ListCachesRequest) (*models.ListCachesResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	resp, err := client.grpcClient.ListCaches(ctx, &pb.XListCachesRequest{NextToken: request.NextToken})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return models.NewListCacheResponse(resp), nil
}

func (client *ScsControlClient) CreateSigningKey(ctx context.Context, endpoint string, request *models.CreateSigningKeyRequest) (*models.CreateSigningKeyResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	resp, err := client.grpcClient.CreateSigningKey(ctx, &pb.XCreateSigningKeyRequest{TtlMinutes: request.TtlMinutes})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	createResp, err := models.NewCreateSigningKeyResponse(endpoint, resp)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return createResp, nil
}

func (client *ScsControlClient) RevokeSigningKey(ctx context.Context, request *models.RevokeSigningKeyRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.RevokeSigningKey(ctx, &pb.XRevokeSigningKeyRequest{KeyId: request.KeyId})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsControlClient) ListSigningKeys(ctx context.Context, endpoint string, request *models.ListSigningKeysRequest) (*models.ListSigningKeysResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	resp, err := client.grpcClient.ListSigningKeys(ctx, &pb.XListSigningKeysRequest{NextToken: request.NextToken})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return models.NewListSigningKeysResponse(endpoint, resp), nil
}
