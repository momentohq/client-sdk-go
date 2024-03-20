package grpcmanagers

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
)

type ScsControlGrpcManager struct {
	Conn *grpc.ClientConn
}

const ControlPort = ":443"

func NewScsControlGrpcManager(request *models.ControlGrpcManagerRequest) (*ScsControlGrpcManager, momentoerrors.MomentoSvcErr) {
	authToken := request.CredentialProvider.GetAuthToken()
	endpoint := fmt.Sprint(request.CredentialProvider.GetControlEndpoint(), ControlPort)

	// Override grpc config to disable keepalives
	controlConfig := config.NewStaticGrpcConfiguration(&config.GrpcConfigurationProps{}).WithKeepAliveDisabled()

	// TODO make NewClient
	conn, err := grpc.Dial(
		endpoint,
		AllDialOptions(
			controlConfig,
			grpc.WithUnaryInterceptor(interceptor.AddAuthHeadersInterceptor(authToken)),
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
