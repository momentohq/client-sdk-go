package grpcmanagers

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
)

type LeaderboardGrpcManager struct {
	Conn *grpc.ClientConn
}

const LeaderboardPort = ":443"

func NewLeaderboardGrpcManager(request *models.LeaderboardGrpcManagerRequest) (*LeaderboardGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), LeaderboardPort)
	authToken := request.CredentialProvider.GetAuthToken()
	// TODO make NewClient
	conn, err := grpc.Dial(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			grpc.WithUnaryInterceptor(interceptor.AddAuthHeadersInterceptor(authToken)),
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
