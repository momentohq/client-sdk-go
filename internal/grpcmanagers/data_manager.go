package grpcmanagers

import (
	"context"
	"errors"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

const CachePort = ":443"

func NewUnaryDataGrpcManager(request *models.DataGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), CachePort)
	authToken := request.CredentialProvider.GetAuthToken()

	var conn *grpc.ClientConn
	var err error
	if request.EagerConnect {
		conn, err = grpc.NewClient(
			endpoint,
			AllDialOptions(
				request.GrpcConfiguration,
				grpc.WithChainUnaryInterceptor(interceptor.AddUnaryRetryInterceptor(request.RetryStrategy), interceptor.AddReadConcernHeaderInterceptor(request.ReadConcern), interceptor.AddAuthHeadersInterceptor(authToken)),
				grpc.WithChainStreamInterceptor(interceptor.AddStreamHeaderInterceptor(authToken)),
			)...,
		)
	} else {
		conn, err = grpc.NewClient(
			endpoint,
			AllDialOptions(
				request.GrpcConfiguration,
				grpc.WithChainUnaryInterceptor(interceptor.AddUnaryRetryInterceptor(request.RetryStrategy), interceptor.AddReadConcernHeaderInterceptor(request.ReadConcern), interceptor.AddAuthHeadersInterceptor(authToken)),
				grpc.WithChainStreamInterceptor(interceptor.AddStreamHeaderInterceptor(authToken)),
			)...,
		)
	}

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
				return nil
			case connectivity.Idle, connectivity.Connecting:
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
