package grpcmanagers

import (
	"crypto/tls"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	internalRequests "github.com/momentohq/client-sdk-go/internal/requests"
	"github.com/momentohq/client-sdk-go/internal/scserrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ControlGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewControlGrpcManager(cgmr internalRequests.ControlGrpcManagerRequest) (*ControlGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(cgmr.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(cgmr.AuthToken)))
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	return &ControlGrpcManager{Conn: conn}, nil
}

func (cm *ControlGrpcManager) Close() error {
	return cm.Conn.Close()
}
