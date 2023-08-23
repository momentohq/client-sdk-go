package grpcmanagers

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

const CachePort = ":443"

func NewUnaryDataGrpcManager(request *models.DataGrpcManagerRequest, eagerConnectTimeout time.Duration) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), CachePort)
	authToken := request.CredentialProvider.GetAuthToken()

	var conn *grpc.ClientConn
	var err error
	if eagerConnectTimeout > 0 {
		conn, err = grpc.Dial(
			endpoint,
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithTransportCredentials(credentials.NewTLS(config)),
			grpc.WithChainUnaryInterceptor(interceptor.AddUnaryRetryInterceptor(request.RetryStrategy), interceptor.AddHeadersInterceptor(authToken)),
		)
	} else {
		conn, err = grpc.Dial(
			endpoint,
			grpc.WithTransportCredentials(credentials.NewTLS(config)),
			grpc.WithChainUnaryInterceptor(interceptor.AddUnaryRetryInterceptor(request.RetryStrategy), interceptor.AddHeadersInterceptor(authToken)),
		)
	}

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}

func NewStreamDataGrpcManager(request *models.DataStreamGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), CachePort)
	authToken := request.CredentialProvider.GetAuthToken()
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithChainStreamInterceptor(interceptor.AddStreamHeaderInterceptor(authToken)),
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

func (gm *DataGrpcManager) Connect(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			// Context timeout or cancellation occurred
			return ctx.Err()
		default:
			// Check current state
			state := gm.Conn.GetState()
			switch state {
			case connectivity.Ready:
				// Connection is ready, exit the method
				fmt.Println("connection is ready")
				return nil
			case connectivity.Idle, connectivity.Connecting:
				fmt.Println("connection is " + state.String())
				// If Idle or Connecting, wait for a state change
				if !gm.Conn.WaitForStateChange(ctx, state) {
					// Context was done while waiting
					return ctx.Err()
				}
			default:
				// For other states like TransientFailure, you might want to handle them differently
				return errors.New("connection is in an unexpected state")
			}
		}
	}
}
