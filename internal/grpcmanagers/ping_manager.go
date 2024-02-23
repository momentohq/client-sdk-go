package grpcmanagers

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
)

type PingGrpcManager struct {
	Conn *grpc.ClientConn
}

const PingPort = ":443"

func NewPingGrpcManager(request *models.PingGrpcManagerRequest) (*PingGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), PingPort)
	authToken := request.CredentialProvider.GetAuthToken()
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
	return &PingGrpcManager{Conn: conn}, nil
}

func (pingManager *PingGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(pingManager.Conn.Close())
}
