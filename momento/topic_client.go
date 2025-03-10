// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type TopicClient interface {
	Subscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error)
	Publish(ctx context.Context, request *TopicPublishRequest) (responses.TopicPublishResponse, error)

	Close()
}

// defaultTopicClient represents all information needed for momento client to enable publish and subscribe operations.
type defaultTopicClient struct {
	credentialProvider auth.CredentialProvider
	numChannels        uint32
	pubSubClient       *pubSubClient
	log                logger.MomentoLogger
}

// NewTopicClient returns a new TopicClient with provided configuration and credential provider arguments.
func NewTopicClient(topicsConfiguration config.TopicsConfiguration, credentialProvider auth.CredentialProvider) (TopicClient, error) {
	numChannels := topicsConfiguration.GetNumGrpcChannels()

	client := &defaultTopicClient{
		credentialProvider: credentialProvider,
		numChannels:        numChannels,
		log:                topicsConfiguration.GetLoggerFactory().GetLogger("topic-client"),
	}

	pubSubClient, err := newPubSubClient(&models.PubSubClientRequest{
		CredentialProvider:  credentialProvider,
		TopicsConfiguration: topicsConfiguration,
		Log:                 client.log,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.pubSubClient = pubSubClient

	return client, nil
}

func (c defaultTopicClient) Subscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.TopicName, "Topic name"); err != nil {
		return nil, err
	}

	var topicManager *grpcmanagers.TopicGrpcManager
	var subscribeClient pb.Pubsub_SubscribeClient
	var cancelContext context.Context
	var cancelFunction context.CancelFunc
	var err error

	var firstMsg *pb.XSubscriptionItem
	topicManager, subscribeClient, cancelContext, cancelFunction, err = c.pubSubClient.topicSubscribe(ctx, &TopicSubscribeRequest{
		CacheName:                   request.CacheName,
		TopicName:                   request.TopicName,
		ResumeAtTopicSequenceNumber: request.ResumeAtTopicSequenceNumber,
		SequencePage:                request.SequencePage,
	})
	if err != nil {
		return nil, err
	}

	if request.ResumeAtTopicSequenceNumber == 0 && request.SequencePage == 0 {
		c.log.Debug("Starting new subscription with new sequence number and sequence page.")
	} else {
		c.log.Debug("Resuming subscription from sequence number %d and sequence page %d.", request.ResumeAtTopicSequenceNumber, request.SequencePage)
	}

	// Ping the stream to provide a nice error message if the cache does not exist.
	firstMsg, err = subscribeClient.Recv()
	if err != nil {
		c.log.Debug("failed to receive first message from subscription: %s", err.Error())

		// We now count number of active subscriptions per grpc channel, so if we did not return
		// an error earlier when calling c.pubSubClient.topicSubscribe, we know that the error
		// here is due to a service-side subscription limit.
		rpcError, _ := status.FromError(err)
		if rpcError != nil {
			if rpcError.Code() == codes.ResourceExhausted {
				c.log.Warn("Topic subscription limit reached for this account; please contact us at support@momentohq.com")
			}
		}
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	switch firstMsg.Kind.(type) {
	case *pb.XSubscriptionItem_Heartbeat:
		// The first message to a new subscription will always be a heartbeat.
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InternalServerError,
			fmt.Sprintf("expected a heartbeat message, got: %T", firstMsg.Kind),
			err,
		)
	}

	return &topicSubscription{
		topicManager:       topicManager,
		subscribeClient:    subscribeClient,
		momentoTopicClient: c.pubSubClient,
		cacheName:          request.CacheName,
		topicName:          request.TopicName,
		log:                c.log,
		cancelContext:      cancelContext,
		cancelFunction:     cancelFunction,
	}, nil
}

func (c defaultTopicClient) Publish(ctx context.Context, request *TopicPublishRequest) (responses.TopicPublishResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.TopicName, "Topic name"); err != nil {
		return nil, err
	}

	if request.Value == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "value cannot be nil", nil,
			),
		)
	}

	err := c.pubSubClient.topicPublish(ctx, &TopicPublishRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
		Value:     request.Value,
	})

	if err != nil {
		c.log.Debug("failed to topic publish: %s", err.Error())
		return nil, err
	}

	return &responses.TopicPublishSuccess{}, err
}

func (c defaultTopicClient) Close() {
	defer c.pubSubClient.close()
}
