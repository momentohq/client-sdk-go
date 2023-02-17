package momento

import (
	"context"
	"io"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TopicSubscription interface {
	Item(ctx context.Context) (TopicValue, error)
}

type topicSubscription struct {
	grpcClient         grpc.ClientStream
	momentoTopicClient *pubSubClient
	cacheName          string
	topicName          string
}

func (s *topicSubscription) Item(ctx context.Context) (TopicValue, error) {
	for {
		rawMsg := new(pb.XSubscriptionItem)
		if err := s.grpcClient.RecvMsg(rawMsg); err != nil {
			err := s.handleStreamError(ctx, err)
			if err != nil {
				return nil, err
			}
		} else {
			switch typedMsg := rawMsg.Kind.(type) {
			case *pb.XSubscriptionItem_Discontinuity:
				continue
			case *pb.XSubscriptionItem_Item:
				switch subscriptionItem := typedMsg.Item.Value.Kind.(type) {
				case *pb.XTopicValue_Text:
					return &TopicValueString{
						Text: subscriptionItem.Text,
					}, nil
				case *pb.XTopicValue_Binary:
					return &TopicValueBytes{
						Bytes: subscriptionItem.Binary,
					}, nil
				}
			case *pb.XSubscriptionItem_Heartbeat:
				// Doesn't count against our retries.
				continue
			default:
				// Ignore unknown responses, so we don't stop polling if we add a new message.
				// For example, we wouldn't want to stop because of an unknown heartbeat response.
				continue
			}
		}
	}
}

func (s *topicSubscription) handleStreamError(ctx context.Context, err error) error {
	var returnErr error
	if err == io.EOF {
		returnErr = s.reInitStream(ctx)
	} else if grpcStatusErr, ok := status.FromError(err); ok {
		if grpcStatusErr.Code() == codes.Internal &&
			grpcStatusErr.Message() == "stream terminated by RST_STREAM with error code: NO_ERROR" {
			returnErr = s.reInitStream(ctx)
		}
	}
	if returnErr != nil {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InternalServerError,
			"Fatal error getting item from topic",
			err,
		)
	}
	return nil
}

func (s *topicSubscription) reInitStream(ctx context.Context) error {
	newStream, err := s.momentoTopicClient.TopicSubscribe(ctx, &TopicSubscribeRequest{
		CacheName: s.cacheName,
		TopicName: s.topicName,
	})
	if err != nil {
		return err
	}
	s.grpcClient = newStream
	return nil
}
