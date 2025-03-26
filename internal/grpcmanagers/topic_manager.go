package grpcmanagers

import (
	"context"
	"fmt"
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
}

func NewStreamTopicGrpcManager(request *models.TopicStreamGrpcManagerRequest) (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetCacheEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	middlewareList := request.Middleware
	var onRequestCallback func(context.Context, string)
	var onStreamRequestCallback func(context.Context, string)
	for _, mw := range middlewareList {
		fmt.Printf("\n=====\nstream topic manager looking at %T middleware\n", mw)
		if rmw, ok := mw.(middleware.InterceptorCallbackMiddleware); ok {
			onRequestCallback = rmw.OnInterceptorRequest
			onStreamRequestCallback = rmw.OnStreamInterceptorRequest
			break
		}
	}

	headerInterceptors := []grpc.StreamClientInterceptor{
		interceptor.AddStreamRetryInterceptor(request.RetryStrategy, onStreamRequestCallback),
		interceptor.AddStreamHeaderInterceptor(authToken),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsCacheEndpointSecure(),
			grpc.WithChainStreamInterceptor(headerInterceptors...),
			grpc.WithChainUnaryInterceptor(
				interceptor.AddAuthHeadersInterceptor(authToken),
				interceptor.AddUnaryRetryInterceptor(request.RetryStrategy, onRequestCallback),
			),
		)...,
	)

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &TopicGrpcManager{
		Conn:         conn,
		StreamClient: pb.NewPubsubClient(conn),
		Middleware:   request.Middleware,
	}, nil
}

func (topicManager *TopicGrpcManager) Close() momentoerrors.MomentoSvcErr {
	topicManager.NumActiveSubscriptions.Store(0)
	return momentoerrors.ConvertSvcErr(topicManager.Conn.Close())
}
