package grpcmanagers

import (
	"crypto/tls"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
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
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(authToken)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor([]grpc_retry.CallOption{
			grpc_retry.WithMax(request.Config.GetRetryStrategy().GetMaxRetries()),
			grpc_retry.WithPerRetryTimeout(request.Config.GetRetryStrategy().GetPerRetryTimeout()),
			grpc_retry.WithCodes(request.Config.GetRetryStrategy().GetRetryableRequestStatuses()...),
		}...)),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &ScsControlGrpcManager{Conn: conn}, nil
}

func (controlManager *ScsControlGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(controlManager.Conn.Close())
}
