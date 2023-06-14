package momento

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type TopicSubscription interface {
	Item(ctx context.Context) (TopicValue, error)
}

type topicSubscription struct {
	topicManager            *grpcmanagers.TopicGrpcManager
	grpcClient              grpc.ClientStream
	momentoTopicClient      *pubSubClient
	cacheName               string
	topicName               string
	log                     logger.MomentoLogger
	lastKnownSequenceNumber uint64
}

func (s *topicSubscription) Item(ctx context.Context) (TopicValue, error) {
	for {
		rawMsg := new(pb.XSubscriptionItem)
		if err := s.grpcClient.RecvMsg(rawMsg); err != nil {
			s.log.Error("stream disconnected, attempting to reconnect err:", fmt.Sprint(err))
			s.attemptReconnect(ctx)
			// retry getting the latest item
			continue
		}

		switch typedMsg := rawMsg.Kind.(type) {
		case *pb.XSubscriptionItem_Discontinuity:
			s.log.Debug("recieved discontinuity item")
			continue
		case *pb.XSubscriptionItem_Item:
			s.lastKnownSequenceNumber = typedMsg.Item.GetTopicSequenceNumber()
			switch subscriptionItem := typedMsg.Item.Value.Kind.(type) {
			case *pb.XTopicValue_Text:
				return String(subscriptionItem.Text), nil
			case *pb.XTopicValue_Binary:
				return Bytes(subscriptionItem.Binary), nil
			}
		case *pb.XSubscriptionItem_Heartbeat:
			s.log.Debug("recieved heartbeat item")
			continue
		default:
			s.log.Trace("Unrecognized response detected.",
				"response", fmt.Sprint(typedMsg))
			continue
		}
	}
}

func (s *topicSubscription) attemptReconnect(ctx context.Context) {
	// make sure the connection isn't ready, if its ready no need to reconnect
	if s.topicManager.Conn.GetState() == connectivity.Ready {
		s.log.Debug("connection is in ready state, not reconnecting")
		return
	}
	// try and reconnect every n seconds. This will attempt to reconnect indefinetly
	seconds := 5 * time.Second
	for {
		s.log.Debug("Attempting reconnecting to client stream")
		time.Sleep(seconds)
		newTopicManager, newStream, err := s.momentoTopicClient.topicSubscribe(ctx, &TopicSubscribeRequest{
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
			return
		}
	}
}
