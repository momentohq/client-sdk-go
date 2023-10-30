package momento

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type pubSubClient struct {
	streamTopicManagers []*grpcmanagers.TopicGrpcManager
	endpoint            string
	log                 logger.MomentoLogger
}

var streamTopicManagerCount uint64
var numGrpcStreams int64
var numChannels uint32

func newPubSubClient(request *models.PubSubClientRequest) (*pubSubClient, momentoerrors.MomentoSvcErr) {
	numSubscriptions := float64(request.TopicsConfiguration.GetMaxSubscriptions())
	numGrpcChannels := request.TopicsConfiguration.GetNumGrpcChannels()

	// numSubscriptions is deprecated. Nevertheless, check that numGrpcChannels and numSubscriptions
	// are not both set. They should be mutually exclusive configs.
	if numGrpcChannels > 0 && numSubscriptions > 0 {
		return nil, NewMomentoError(momentoerrors.InvalidArgumentError, "Cannot accept both maxSubscriptions and numGrpcChannels as arguments; please use numGrpcChannels as maxSubscriptions is deprecated", nil)
	}

	if numGrpcChannels > 0 {
		numChannels = numGrpcChannels
	} else if numSubscriptions > 0 {
		// a single channel can support 100 streams, so we need to create enough
		// channels to handle the maximum number of subscriptions
		// plus one for the publishing channel
		numChannels = uint32(math.Ceil((numSubscriptions + 1) / 100.0))
	} else {
		numChannels = 1
	}
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)

	for i := 0; uint32(i) < numChannels; i++ {
		streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(&models.TopicStreamGrpcManagerRequest{
			CredentialProvider: request.CredentialProvider,
		})
		if err != nil {
			return nil, err
		}
		streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	}

	return &pubSubClient{
		streamTopicManagers: streamTopicManagers,
		endpoint:            request.CredentialProvider.GetCacheEndpoint(),
		log:                 request.Log,
	}, nil
}

func (client *pubSubClient) getNextStreamTopicManager() *grpcmanagers.TopicGrpcManager {
	nextManagerIndex := atomic.AddUint64(&streamTopicManagerCount, 1)
	topicManager := client.streamTopicManagers[nextManagerIndex%uint64(len(client.streamTopicManagers))]
	return topicManager
}

func (client *pubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, grpc.ClientStream, error) {

	checkNumConcurrentStreams(client.log)

	atomic.AddInt64(&numGrpcStreams, 1)

	topicManager := client.getNextStreamTopicManager()
	clientStream, err := topicManager.StreamClient.Subscribe(ctx, &pb.XSubscriptionRequest{
		CacheName:                   request.CacheName,
		Topic:                       request.TopicName,
		ResumeAtTopicSequenceNumber: request.ResumeAtTopicSequenceNumber,
	})

	if err != nil {
		atomic.AddInt64(&numGrpcStreams, -1)
		return nil, nil, err
	}

	if numGrpcStreams > 0 && (int64(numChannels*100)-numGrpcStreams < 10) {
		fmt.Printf("WARNING: approaching grpc maximum concurrent stream limit, %d remaining of total %d streams\n", int64(numChannels*100)-numGrpcStreams, numChannels*100)
	}

	return topicManager, clientStream, err
}

func (client *pubSubClient) topicPublish(ctx context.Context, request *TopicPublishRequest) error {
	checkNumConcurrentStreams(client.log)

	atomic.AddInt64(&numGrpcStreams, 1)

	topicManager := client.getNextStreamTopicManager()
	switch value := request.Value.(type) {
	case String:
		_, err := topicManager.StreamClient.Publish(ctx, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Text{
					Text: value.asString(),
				},
			},
		})
		atomic.AddInt64(&numGrpcStreams, -1)
		return err
	case Bytes:
		_, err := topicManager.StreamClient.Publish(ctx, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Binary{
					Binary: value.asBytes(),
				},
			},
		})
		atomic.AddInt64(&numGrpcStreams, -1)
		return err
	default:
		atomic.AddInt64(&numGrpcStreams, -1)
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"error encoding topic value only support []byte or string currently", nil,
		)
	}
}

func (client *pubSubClient) close() {
	atomic.AddInt64(&numGrpcStreams, -numGrpcStreams)
	for clientIndex := range client.streamTopicManagers {
		defer client.streamTopicManagers[clientIndex].Close()
	}
}

func checkNumConcurrentStreams(log logger.MomentoLogger) {
	if numGrpcStreams > 0 && numGrpcStreams >= int64(numChannels*100) {
		log.Debug("Already at maximum number of concurrent grpc streams, cannot make new publish or subscribe requests")
	}
}
