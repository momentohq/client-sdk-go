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

type ScsControlGrpcManager struct {
	Conn *grpc.ClientConn
}

const ControlPort = ":443"

func NewScsControlGrpcManager(request *models.ControlGrpcManagerRequest) (*ScsControlGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	authToken := request.CredentialProvider.GetAuthToken()
	endpoint := fmt.Sprint(request.CredentialProvider.GetControlEndpoint(), ControlPort)
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(authToken)),
		grpc.WithUnaryInterceptor(interceptor.AddUnaryRetryInterceptor(request.RetryStrategy)),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &ScsControlGrpcManager{Conn: conn}, nil
}

func (controlManager *ScsControlGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(controlManager.Conn.Close())
}
