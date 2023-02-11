package momento

import (
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc/metadata"
)

const defaultRequestTimeout = 5 * time.Second

type scsDataClient struct {
	grpcManager    *grpcmanagers.DataGrpcManager
	grpcClient     pb.ScsClient
	defaultTtl     time.Duration
	requestTimeout time.Duration
	endpoint       string
}

func newScsDataClient(request *models.DataClientRequest) (*scsDataClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewUnaryDataGrpcManager(&models.DataGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
	})
	if err != nil {
		return nil, err
	}
	var timeout time.Duration
	if request.Configuration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.Configuration.GetClientSideTimeout()
	}
	return &scsDataClient{
		grpcManager:    dataManager,
		grpcClient:     pb.NewScsClient(dataManager.Conn),
		defaultTtl:     request.DefaultTtl,
		requestTimeout: timeout,
		endpoint:       request.CredentialProvider.GetCacheEndpoint(),
	}, nil
}

func (client scsDataClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (scsDataClient) CreateNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
