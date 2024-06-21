package grpcmanagers

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"google.golang.org/grpc"
)

type StoreGrpcManager struct {
	Conn *grpc.ClientConn
}

const StoragePort = ":443"

func NewStoreGrpcManager(request *models.StoreGrpcManagerRequest) (*StoreGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := fmt.Sprint(request.CredentialProvider.GetStorageEndpoint(), StoragePort)
	authToken := request.CredentialProvider.GetAuthToken()
	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			grpc.WithChainUnaryInterceptor(interceptor.AddAuthHeadersInterceptor(authToken)),
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
