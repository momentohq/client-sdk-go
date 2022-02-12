package grpcmanagers

import (
	"crypto/tls"

	"github.com/momentohq/client-sdk-go/internal/models"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ScsDataGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewScsDataGrpcManager(request *models.DataGrpcManagerRequest) (*ScsDataGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(request.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(request.AuthToken)))
	if err != nil {
		return nil, momentoerrors.ConvertError(err)
	}
	return &ScsDataGrpcManager{Conn: conn}, nil
}

func (dataManager *ScsDataGrpcManager) Close() error {
	return dataManager.Conn.Close()
}
