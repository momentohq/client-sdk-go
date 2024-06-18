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

	headerInterceptors := []grpc.UnaryClientInterceptor{
		interceptor.AddAuthHeadersInterceptor(authToken),
	}

	if !interceptor.FirstTimeHeadersSent {
		interceptor.FirstTimeHeadersSent = true
		headerInterceptors = append(headerInterceptors, interceptor.AddRuntimeVersionHeaderInterceptor())
		headerInterceptors = append(headerInterceptors, interceptor.AddAgentHeaderInterceptor("ping"))
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			grpc.WithChainUnaryInterceptor(headerInterceptors...),
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
