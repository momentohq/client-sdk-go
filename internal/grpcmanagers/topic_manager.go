package grpcmanagers

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
)

type TopicGrpcManager struct {
	Conn         *grpc.ClientConn
	StreamClient pb.PubsubClient
}

func NewStreamTopicGrpcManager(request *models.TopicStreamGrpcManagerRequest) (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), CachePort)
	authToken := request.CredentialProvider.GetAuthToken()

	headerInterceptors := []grpc.StreamClientInterceptor{
		interceptor.AddStreamHeaderInterceptor(authToken),
	}

	if !interceptor.FirstTimeHeadersSent {
		interceptor.FirstTimeHeadersSent = true
		headerInterceptors = append(headerInterceptors, interceptor.AddStreamRuntimeVersionHeaderInterceptor())
		headerInterceptors = append(headerInterceptors, interceptor.AddStreamAgentHeaderInterceptor("topic"))
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
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
	}, nil
}

func (topicManager *TopicGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(topicManager.Conn.Close())
}
