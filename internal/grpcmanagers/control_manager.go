package grpcmanagers

import (
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
)

type ScsControlGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewScsControlGrpcManager(request *models.ControlGrpcManagerRequest) (*ScsControlGrpcManager, momentoerrors.MomentoSvcErr) {
	authToken := request.CredentialProvider.GetAuthToken()
	endpoint := request.CredentialProvider.GetControlEndpoint()

	// Override grpc config to disable keepalives
	controlConfig := config.NewStaticGrpcConfiguration(&config.GrpcConfigurationProps{}).WithKeepAliveDisabled()

	headerInterceptors := []grpc.UnaryClientInterceptor{
		interceptor.AddAuthHeadersInterceptor(authToken),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			controlConfig,
			request.CredentialProvider.IsControlEndpointSecure(),
			grpc.WithChainUnaryInterceptor(headerInterceptors...),
		)...,
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &ScsControlGrpcManager{Conn: conn}, nil
}

func (controlManager *ScsControlGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(controlManager.Conn.Close())
}
