package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

type shardedTopicClient struct {
	topicClient       momento.TopicClient
	numShardsPerTopic int
	log               logger.MomentoLogger
}

func NewShardedTopicClient(topicsConfiguration config.TopicsConfiguration, credentialProvider auth.CredentialProvider, numShardsPerTopic int) (momento.TopicClient, error) {
	topicClient, err := momento.NewTopicClient(topicsConfiguration, credentialProvider)
	if err != nil {
		return nil, err
	}
	return shardedTopicClient{
		topicClient:       topicClient,
		numShardsPerTopic: numShardsPerTopic,
		log:               topicsConfiguration.GetLoggerFactory().GetLogger("sharded-topic-client"),
	}, nil
}

func (c shardedTopicClient) Subscribe(ctx context.Context, request *momento.TopicSubscribeRequest) (momento.TopicSubscription, error) {
	subscriptions := make([]namedTopicSubscription, c.numShardsPerTopic)
	for i := 0; i < c.numShardsPerTopic; i++ {
		topicName := fmt.Sprintf("%s-%d", request.TopicName, i)
		subscription, err := c.topicClient.Subscribe(ctx, &momento.TopicSubscribeRequest{
			CacheName: request.CacheName,
			TopicName: topicName,
		})
		if err != nil {
			for j := 0; j < c.numShardsPerTopic; j++ {
				if subscriptions[j].subscription != nil {
					subscriptions[j].subscription.Close()
				}
			}
			return nil, err
		}
		subscriptions[i] = namedTopicSubscription{
			topicName:    topicName,
			subscription: subscription,
		}
	}
	return newShardedTopicSubscription(ctx, subscriptions, c.log), nil
}

func (c shardedTopicClient) Publish(ctx context.Context, request *momento.TopicPublishRequest) (responses.TopicPublishResponse, error) {
	topicNamePrefix := request.TopicName
	shardToPublishTo := rand.Intn(c.numShardsPerTopic)
	topicName := fmt.Sprintf("%s-%d", topicNamePrefix, shardToPublishTo)
	c.log.Debug("Publishing to topic %s", topicName)
	shardRequest := momento.TopicPublishRequest{
		CacheName: cacheName,
		TopicName: topicName,
		Value:     request.Value,
	}
	return c.topicClient.Publish(ctx, &shardRequest)
}

func (c shardedTopicClient) Close() {
	c.topicClient.Close()
}

type namedTopicSubscription struct {
	topicName    string
	subscription momento.TopicSubscription
}

type shardedTopicItem struct {
	value momento.TopicValue
	err   error
}

type shardedTopicSubscription struct {
	cancelContext        context.Context
	cancelFunction       context.CancelFunc
	receivedItemsChannel chan shardedTopicItem
	subscriptions        []namedTopicSubscription
	wg                   *sync.WaitGroup
	log                  logger.MomentoLogger
}

func newShardedTopicSubscription(
	ctx context.Context, subscriptions []namedTopicSubscription, log logger.MomentoLogger,
) shardedTopicSubscription {
	// We want the channel to support a bit of buffer in case we are reading things from the
	// sharded subscriptions faster than the caller is consuming them. Arbitrarily choosing a value
	// for now that should give us plenty of breathing room but not consume an excessive amount of memory.
	itemsChannelSize := 1_000
	receivedItemsChannel := make(chan shardedTopicItem, itemsChannelSize)

	// try withCancel on context
	cancelContext, cancelFunction := context.WithCancel(ctx)

	var wg sync.WaitGroup

	for i := 0; i < len(subscriptions); i++ {
		subscription := subscriptions[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			publishSubscriptionItemsToChannel(ctx, cancelContext, subscription, receivedItemsChannel, log)
			log.Debug("Back from publishSubscriptionItemsToChannel for topic: %s", subscription.topicName)
		}()
	}

	return shardedTopicSubscription{
		cancelContext:        cancelContext,
		cancelFunction:       cancelFunction,
		receivedItemsChannel: receivedItemsChannel,
		subscriptions:        subscriptions,
		wg:                   &wg,
		log:                  log,
	}
}

func (s shardedTopicSubscription) Item(ctx context.Context) (momento.TopicValue, error) {

	for {
		// Its totally possible a client just calls `cancel` on the `context` immediately after subscribing to an
		// item, so we should check that here.
		select {
		case <-ctx.Done():
			{
				s.log.Debug("Context is Done, sharded subscription exiting item loop")
			}
			return nil, ctx.Err()
		case <-s.cancelContext.Done():
			s.log.Debug("Context is Cancelled, sharded subscription exiting item loop")
			return nil, s.cancelContext.Err()

		default:
			// Proceed as is
		}

		item := <-s.receivedItemsChannel
		return item.value, item.err
	}
}

func publishSubscriptionItemsToChannel(
	ctx context.Context, cancelContext context.Context, sub namedTopicSubscription, topicValueChan chan shardedTopicItem, log logger.MomentoLogger) {
	topicName := sub.topicName
	log.Debug("Beginning publish coroutine for topic: %s", topicName)
	for {
		log.Debug("Next iteration of publish coroutine for topic: %s", topicName)
		select {
		case <-ctx.Done():
			// Context has been canceled, return an error
			{
				log.Debug("Context is Done, stopping publish coroutine for topic: %s", topicName)
				return
			}
		case <-cancelContext.Done():
			{
				log.Debug("Context has been cancelled, stopping publish coroutine for topic: %s", topicName)
				return
			}
		default:
			// Proceed as is
		}

		log.Debug("Waiting for next subscription item in coroutine for topic: %s", topicName)
		item, err := sub.subscription.Item(ctx)
		log.Debug("subscription.Item returned a value in coroutine for topic: %s", topicName)
		topicValueChan <- shardedTopicItem{
			value: item,
			err:   err,
		}
	}
}

func (s shardedTopicSubscription) Close() {
	for i := 0; i < len(s.subscriptions); i++ {
		s.log.Debug("Closing subscription for topic: %s", s.subscriptions[i].topicName)
		s.subscriptions[i].subscription.Close()
	}
	s.cancelFunction()
	s.log.Debug("Waiting for sharded subscription goroutines to exit.")
	s.wg.Wait()
	s.log.Debug("All sharded subscription goroutines have exited; ShardedTopicSubscription.Close complete")
}
