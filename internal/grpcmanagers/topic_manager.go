package grpcmanagers

import (
	"crypto/tls"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TopicGrpcManager struct {
	Conn         *grpc.ClientConn
	StreamClient pb.PubsubClient
}

func NewStreamTopicGrpcManager(request *models.TopicStreamGrpcManagerRequest) (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	endpoint := fmt.Sprint(request.CredentialProvider.GetCacheEndpoint(), CachePort)
	authToken := request.CredentialProvider.GetAuthToken()
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithChainStreamInterceptor(interceptor.AddStreamHeaderInterceptor(authToken)),
		grpc.WithChainUnaryInterceptor(interceptor.AddAuthHeadersInterceptor(authToken)),
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
