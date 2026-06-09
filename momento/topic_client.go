// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/config/retry"

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
	// Subscribe starts a subscription on the given topic.
	//
	// The ctx bounds the subscription's lifetime, not just the Subscribe call:
	// the underlying stream — including any stream re-established by the
	// automatic reconnect logic — is a child of this ctx, so cancelling it
	// ends the subscription and releases its connection-pool slot.
	Subscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error)
	Publish(ctx context.Context, request *TopicPublishRequest) (responses.TopicPublishResponse, error)

	Close()
}

// defaultTopicClient represents all information needed for momento client to enable publish and subscribe operations.
type defaultTopicClient struct {
	credentialProvider auth.CredentialProvider
	pubSubClient       *pubSubClient
	log                logger.MomentoLogger
	requestTimeout     time.Duration
	retryStrategy      retry.Strategy
}

// NewTopicClient returns a new TopicClient with provided configuration and credential provider arguments.
func NewTopicClient(topicsConfiguration config.TopicsConfiguration, credentialProvider auth.CredentialProvider) (TopicClient, error) {
	var timeout time.Duration
	if topicsConfiguration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = topicsConfiguration.GetClientSideTimeout()
	}

	client := &defaultTopicClient{
		credentialProvider: credentialProvider,
		log:                topicsConfiguration.GetLoggerFactory().GetLogger("topic-client"),
		requestTimeout:     timeout,
		retryStrategy:      topicsConfiguration.GetRetryStrategy(),
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

	// Set a timeout by which the first heartbeat message should be received.
	// If the first message is not received within this time, we will cancel the subscription.
	firstMessageCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()
	subChan := make(chan *topicSubscription, 1)
	errChan := make(chan error, 1)

	go c.sendSubscribe(ctx, request, subChan, errChan)
	select {
	case <-ctx.Done():
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.CanceledError,
			"subscribe request context was canceled",
			ctx.Err(),
		)
	case <-firstMessageCtx.Done():
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.TimeoutError,
			"subscription did not receive first message within the expected time",
			nil,
		)
	case subscription := <-subChan:
		return subscription, nil
	case err := <-errChan:
		return nil, err
	}
}

func (c defaultTopicClient) sendSubscribe(requestCtx context.Context, request *TopicSubscribeRequest, subChan chan *topicSubscription, errChan chan error) {
	state, err := c.pubSubClient.topicSubscribe(requestCtx, &TopicSubscribeRequest{
		CacheName:                   request.CacheName,
		TopicName:                   request.TopicName,
		ResumeAtTopicSequenceNumber: request.ResumeAtTopicSequenceNumber,
		SequencePage:                request.SequencePage,
	})
	if err != nil {
		errChan <- err
		return
	}

	if request.ResumeAtTopicSequenceNumber == 0 && request.SequencePage == 0 {
		c.log.Debug("Starting new subscription with new sequence number and sequence page.")
	} else {
		c.log.Debug("Resuming subscription from sequence number %d and sequence page %d.", request.ResumeAtTopicSequenceNumber, request.SequencePage)
	}

	// Single cleanup path: cancel, release, report.
	fail := func(err error) {
		state.cancelFunction()
		state.reservation.Release()
		errChan <- err
	}

	// Ping the stream to surface a nice error if the cache does not exist.
	firstMsg, err := state.subscribeClient.Recv()
	if err != nil {
		c.log.Debug("failed to receive first message from subscription: %s", err.Error())

		// topicSubscribe would have errored already if the local pool was
		// exhausted, so a ResourceExhausted here is a service-side limit.
		rpcError, _ := status.FromError(err)
		if rpcError != nil {
			if rpcError.Code() == codes.ResourceExhausted {
				c.log.Warn("Topic subscription limit reached for this account; please contact us at support@momentohq.com")
			}
		}
		fail(momentoerrors.ConvertSvcErr(err))
		return
	}

	switch firstMsg.Kind.(type) {
	case *pb.XSubscriptionItem_Heartbeat:
		// The first message to a new subscription will always be a heartbeat.
	default:
		fail(momentoerrors.NewMomentoSvcErr(
			momentoerrors.InternalServerError,
			fmt.Sprintf("expected a heartbeat message, got: %T", firstMsg.Kind),
			err,
		))
		return
	}

	var topicEventCallback func(cacheName string, requestName string, event middleware.TopicSubscriptionEventType)
	for _, mw := range c.pubSubClient.middleware {
		if rmw, ok := mw.(middleware.TopicEventCallbackMiddleware); ok {
			// currently this is exclusively used for resubscribe metrics in the MomentoLocalMiddleware,
			// so we break after we find the first one.
			topicEventCallback = rmw.OnTopicEvent
			break
		}
	}
	sub := &topicSubscription{
		topicEventCallback: topicEventCallback,
		momentoTopicClient: c.pubSubClient,
		cacheName:          request.CacheName,
		topicName:          request.TopicName,
		log:                c.log,
		retryStrategy:      c.retryStrategy,
		streamParentCtx:    requestCtx,
		closedSignal:       make(chan struct{}),
	}
	sub.state.Store(state)
	subChan <- sub
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
