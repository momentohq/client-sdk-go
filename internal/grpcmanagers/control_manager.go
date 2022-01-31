package grpcmanagers

import (
	"crypto/tls"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ScsControlGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewScsControlGrpcManager(request *models.ControlGrpcManagerRequest) (*ScsControlGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(request.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(request.AuthToken)))
	if err != nil {
		return nil, momentoerrors.GrpcErrorConverter(err)
	}
	return &ScsControlGrpcManager{Conn: conn}, nil
}

func (controlManager *ScsControlGrpcManager) Close() error {
	return controlManager.Conn.Close()
}
