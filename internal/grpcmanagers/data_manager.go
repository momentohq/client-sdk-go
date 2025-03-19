package grpcmanagers

import (
	"context"
	"errors"

	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewUnaryDataGrpcManager(request *models.DataGrpcManagerRequest) (*DataGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetCacheEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	// Check the middleware list for a "RetryMiddleware" and use the AddUnaryRetryInterceptor method from the
	// first one found. This allows us to swap out the production retry interceptor for the custom retry interceptor
	// in RetryMetricsMiddleware that tracks retry timestamps for testing.
	middlewareList := request.Middleware
	addRetryInterceptor := interceptor.AddUnaryRetryInterceptor
	for _, mw := range middlewareList {
		if rmw, ok := mw.(middleware.RetryMiddleware); ok {
			addRetryInterceptor = rmw.AddUnaryRetryInterceptor
			break
		}
	}

	headerInterceptors := []grpc.UnaryClientInterceptor{
		addRetryInterceptor(request.RetryStrategy),
		interceptor.AddReadConcernHeaderInterceptor(request.ReadConcern),
		interceptor.AddAuthHeadersInterceptor(authToken),
	}

	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsCacheEndpointSecure(),
			grpc.WithChainUnaryInterceptor(headerInterceptors...),
			grpc.WithChainStreamInterceptor(interceptor.AddStreamHeaderInterceptor(authToken)),
		)...,
	)

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &DataGrpcManager{Conn: conn}, nil
}

func (dataManager *DataGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(dataManager.Conn.Close())
}

func (dataManager *DataGrpcManager) Connect(ctx context.Context) error {
	// grpc.NewClient remains in IDLE until first RPC, but we force
	// an eager connection when we call Connect here
	dataManager.Conn.Connect()

	for {
		select {
		case <-ctx.Done():
			// Context timeout or cancellation occurred
			return ctx.Err()
		default:
			// Check current state
			state := dataManager.Conn.GetState()
			switch state {
			case connectivity.Ready:
				// Connection is ready, exit the method
				return nil
			case connectivity.Idle, connectivity.Connecting:
				// If Idle or Connecting, wait for a state change
				if !dataManager.Conn.WaitForStateChange(ctx, state) {
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
