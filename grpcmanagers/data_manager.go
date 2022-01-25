package grpcmanager

import (
	"crypto/tls"

	interceptor "github.com/momentohq/client-sdk-go/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewDataGrpcManager(authToken string, endPoint string) (DataGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(authToken)))
	return DataGrpcManager{Conn: conn}, err
}

func (cm *DataGrpcManager) Close() error {
	return cm.Conn.Close()
}
