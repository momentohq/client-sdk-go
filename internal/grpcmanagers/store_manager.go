package grpcmanagers

import (
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"google.golang.org/grpc"
)

type StoreGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewStoreGrpcManager(request *models.StoreGrpcManagerRequest) (*StoreGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetStorageEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	headerInterceptors := []grpc.UnaryClientInterceptor{
		interceptor.AddAuthHeadersInterceptor(authToken),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsStorageEndpointSecure(),
			request.CredentialProvider,
			grpc.WithChainUnaryInterceptor(headerInterceptors...),
		)...,
	)

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &StoreGrpcManager{Conn: conn}, nil
}

func (grpcManager *StoreGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(grpcManager.Conn.Close())
}
