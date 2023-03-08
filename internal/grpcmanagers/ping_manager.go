package grpcmanagers

import (
	"crypto/tls"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type PingGrpcManager struct {
	Conn *grpc.ClientConn
}

const PingPort = ":443"

func NewPingGrpcManager(request *models.PingGrpcManagerRequest) (*PingGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), PingPort)
	authToken := request.CredentialProvider.GetAuthToken()
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(authToken)),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &PingGrpcManager{Conn: conn}, nil
}

func (pingManager *PingGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(pingManager.Conn.Close())
}
