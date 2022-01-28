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

func NewDataGrpcManager(request *internalRequests.DataGrpcManagerRequest) (*DataGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(request.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(request.AuthToken)))
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}

func (dataManager *DataGrpcManager) Close() error {
	return dataManager.Conn.Close()
}
