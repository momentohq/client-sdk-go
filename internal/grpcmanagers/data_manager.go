package grpcmanagers

import (
	"crypto/tls"
	"fmt"
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

const CachePort = ":443"

func NewUnaryDataGrpcManager(request *models.DataGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), CachePort)
	authToken := request.CredentialProvider.GetAuthToken()
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithUnaryInterceptor(interceptor.AddHeadersInterceptor(authToken)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor([]grpc_retry.CallOption{
			grpc_retry.WithMax(request.Config.GetRetryStrategy().GetMaxRetries()),
			grpc_retry.WithPerRetryTimeout(request.Config.GetRetryStrategy().GetPerRetryTimeout()),
			grpc_retry.WithCodes(request.Config.GetRetryStrategy().GetRetryableRequestStatuses()...),
		}...)),
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
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), CachePort)
	authToken := request.CredentialProvider.GetAuthToken()
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithDisableRetry(),
		grpc.WithStreamInterceptor(interceptor.AddStreamHeaderInterceptor(authToken)),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor([]grpc_retry.CallOption{
			grpc_retry.WithMax(request.Config.GetRetryStrategy().GetMaxRetries()),
			grpc_retry.WithPerRetryTimeout(request.Config.GetRetryStrategy().GetPerRetryTimeout()),
			grpc_retry.WithCodes(request.Config.GetRetryStrategy().GetRetryableRequestStatuses()...),
		}...)),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}

func NewLocalDataGrpcManager(request *models.LocalDataGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
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
