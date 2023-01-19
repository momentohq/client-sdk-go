package grpcmanagers

import (
	"crypto/tls"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewUnaryDataGrpcManager(request *models.DataGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(
		request.Endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithDisableRetry(),
		grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(request.AuthToken)),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}

func NewStreamDataGrpcManager(request *models.DataGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(
		request.Endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithDisableRetry(),
		grpc.WithStreamInterceptor(interceptor.AddStreamHeaderInterceptor(request.AuthToken)),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}

func NewLocalDataGrpcManager(request *models.DataGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	conn, err := grpc.Dial(
		request.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDisableRetry(),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}
func (dataManager *DataGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(dataManager.Conn.Close())
}
