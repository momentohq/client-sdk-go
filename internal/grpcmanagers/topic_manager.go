package grpcmanagers

import (
	"context"
	"github.com/momentohq/client-sdk-go/config/retry"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config/middleware"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
)

type TopicGrpcManager struct {
	Conn                   *grpc.ClientConn
	StreamClient           pb.PubsubClient
	NumActiveSubscriptions atomic.Int64
	Middleware             []middleware.Middleware
	RetryStrategy retry.Strategy
}

func NewStreamTopicGrpcManager(request *models.TopicStreamGrpcManagerRequest) (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetCacheEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	middlewareList := request.Middleware
	// TODO: this will be called in topic_subscription.go and needs to be added to the TopicGrpcManager struct.
	//  It probably needs to be extended to support callbacks for Item, Error, and Retry events.
	//  These callbacks will probably use a channel to talk to the middleware like the cache client does.
	var onStreamItemCallback func(context.Context, string)
	for _, mw := range middlewareList {
		// TODO: we'll need a different interface for stream middleware
		if rmw, ok := mw.(middleware.InterceptorCallbackMiddleware); ok {
			onStreamItemCallback = rmw.OnStreamInterceptorRequest
			break
		}
	}

	headerInterceptors := []grpc.StreamClientInterceptor{
		interceptor.AddStreamHeaderInterceptor(authToken),
		interceptor.AddStreamRetryInterceptor(),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsCacheEndpointSecure(),
			grpc.WithChainStreamInterceptor(headerInterceptors...),
			grpc.WithChainUnaryInterceptor(interceptor.AddAuthHeadersInterceptor(authToken)),
		)...,
	)

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &TopicGrpcManager{
		Conn:         conn,
		StreamClient: pb.NewPubsubClient(conn),
		Middleware:   request.Middleware,
		RetryStrategy: request.RetryStrategy,
	}, nil
}

func (topicManager *TopicGrpcManager) Close() momentoerrors.MomentoSvcErr {
	topicManager.NumActiveSubscriptions.Store(0)
	return momentoerrors.ConvertSvcErr(topicManager.Conn.Close())
}
