package grpcmanagers

import (
	"crypto/tls"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/scserrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ControlGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewControlGrpcManager(request *models.ControlGrpcManagerRequest) (*ControlGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(request.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(request.AuthToken)))
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	return &ControlGrpcManager{Conn: conn}, nil
}

func (controlManager *ControlGrpcManager) Close() error {
	return controlManager.Conn.Close()
}
