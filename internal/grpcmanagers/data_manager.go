package grpcmanagers

import (
	"crypto/tls"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	internalRequests "github.com/momentohq/client-sdk-go/internal/requests"
	"github.com/momentohq/client-sdk-go/internal/scserrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewDataGrpcManager(dgmr internalRequests.DataGrpcManagerRequest) (*DataGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(dgmr.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(dgmr.AuthToken)))
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}

func (cm *DataGrpcManager) Close() error {
	return cm.Conn.Close()
}
