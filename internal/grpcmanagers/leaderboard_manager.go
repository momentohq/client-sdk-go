package grpcmanagers

import (
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
)

type LeaderboardGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewLeaderboardGrpcManager(request *models.LeaderboardGrpcManagerRequest) (*LeaderboardGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetCacheEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	headerInterceptors := []grpc.UnaryClientInterceptor{
		interceptor.AddAuthHeadersInterceptor(authToken),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsCacheEndpointSecure(),
			grpc.WithChainUnaryInterceptor(headerInterceptors...),
		)...,
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &LeaderboardGrpcManager{Conn: conn}, nil
}

func (grpcManager *LeaderboardGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(grpcManager.Conn.Close())
}
