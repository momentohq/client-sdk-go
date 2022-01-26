package grpcmanagers

import (
	"crypto/tls"

	"github.com/momentohq/client-sdk-go/interceptor"
	"github.com/momentohq/client-sdk-go/scserrors"

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
	return DataGrpcManager{Conn: conn}, scserrors.GrpcErrorConverter(err)
}

func (cm *DataGrpcManager) Close() error {
	return cm.Conn.Close()
}
