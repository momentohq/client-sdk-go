package grpcmanager

import (
	"crypto/tls"

	interceptor "github.com/momentohq/client-sdk-go/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ControlGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewControlGrpcManager(authToken string, endPoint string) (ControlGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(authToken)))
	return ControlGrpcManager{Conn: conn}, err
}

func (cm *ControlGrpcManager) Close() error {
	return cm.Conn.Close()
}
