package momento

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type TopicSubscription interface {
	Item(ctx context.Context) (TopicValue, error)
	DetailedItem(ctx context.Context) (DetailedTopicItem, error)
	Close()
}

type topicSubscription struct {
	topicManager            *grpcmanagers.TopicGrpcManager
	grpcClient              grpc.ClientStream
	momentoTopicClient      *pubSubClient
	cacheName               string
	topicName               string
	log                     logger.MomentoLogger
	lastKnownSequenceNumber uint64
	cancelContext           context.Context
	cancelFunction          context.CancelFunc
}

func (s *topicSubscription) Item(ctx context.Context) (TopicValue, error) {
	for {
		item, err := s.DetailedItem(ctx)
		if err != nil {
			return nil, err
		}

		switch item := item.(type) {
		case TopicItem:
			return item.GetValue(), nil
		case TopicHeartbeat:
			continue
		case TopicDiscontinuity:
			continue
		}
	}
}

func (s *topicSubscription) DetailedItem(ctx context.Context) (DetailedTopicItem, error) {
	for {
		// Its totally possible a client just calls `cancel` on the `context` immediately after subscribing to an
		// item, so we should check that here.
		select {
		case <-ctx.Done():
			// Context has been canceled, return an error
			return nil, ctx.Err()
		case <-s.cancelContext.Done():
			// Context has been canceled, return an error
			return nil, s.cancelContext.Err()
		default:
			// Proceed as is
		}

		rawMsg := new(pb.XSubscriptionItem)
		if err := s.grpcClient.RecvMsg(rawMsg); err != nil {
			select {
			case <-ctx.Done():
				{
					s.log.Info("Subscription context is done; closing subscription.")
					return nil, ctx.Err()
				}
			case <-s.cancelContext.Done():
				{
					s.log.Info("Subscription context is cancelled; closing subscription.")
					return nil, s.cancelContext.Err()
				}
			default:
				{
					// Attempt to reconnect
					s.log.Error("stream disconnected YO, attempting to reconnect err:", fmt.Sprint(err))
					s.attemptReconnect(ctx)
				}
			}

			// retry getting the latest item
			continue
		}

		switch typedMsg := rawMsg.Kind.(type) {
		case *pb.XSubscriptionItem_Discontinuity:
			s.log.Debug("received discontinuity item")
			return NewTopicDiscontinuity(typedMsg.Discontinuity.LastTopicSequence, typedMsg.Discontinuity.NewTopicSequence), nil
		case *pb.XSubscriptionItem_Item:
			s.lastKnownSequenceNumber = typedMsg.Item.GetTopicSequenceNumber()
			switch subscriptionItem := typedMsg.Item.Value.Kind.(type) {
			case *pb.XTopicValue_Text:
				return NewTopicItem(String(subscriptionItem.Text), String(typedMsg.Item.PublisherId), typedMsg.Item.TopicSequenceNumber), nil
			case *pb.XTopicValue_Binary:
				return NewTopicItem(Bytes(subscriptionItem.Binary), String(typedMsg.Item.PublisherId), typedMsg.Item.TopicSequenceNumber), nil
			}
		case *pb.XSubscriptionItem_Heartbeat:
			s.log.Debug("received heartbeat item")
			return TopicHeartbeat{}, nil
		default:
			s.log.Trace("Unrecognized response detected.",
				"response", fmt.Sprint(typedMsg))
			continue
		}
	}
}

func (s *topicSubscription) attemptReconnect(ctx context.Context) {
	// try and reconnect every n seconds. This will attempt to reconnect indefinetly
	seconds := 5 * time.Second
	for {
		s.log.Debug("Attempting reconnecting to client stream")
		time.Sleep(seconds)
		newTopicManager, newStream, cancelContext, cancelFunction, err := s.momentoTopicClient.topicSubscribe(ctx, &TopicSubscribeRequest{
			CacheName:                   s.cacheName,
			TopicName:                   s.topicName,
			ResumeAtTopicSequenceNumber: s.lastKnownSequenceNumber,
		})

		if err != nil {
			s.log.Debug("failed to reconnect to stream, will continue to try in %s seconds", fmt.Sprint(seconds))
		} else {
			s.log.Debug("successfully reconnected to subscription stream")
			s.topicManager = newTopicManager
			s.grpcClient = newStream
			s.cancelContext = cancelContext
			s.cancelFunction = cancelFunction
			return
		}
	}
}

func (s *topicSubscription) Close() {
	atomic.AddInt64(&numGrpcStreams, -1)
	s.cancelFunction()
}
