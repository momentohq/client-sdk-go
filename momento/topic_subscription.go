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
					return String(subscriptionItem.Text), nil
				case *pb.XTopicValue_Binary:
					return Bytes(subscriptionItem.Binary), nil
				}
			case *pb.XSubscriptionItem_Heartbeat:
				// FIXME add warning logging here
				continue
			default:
				// FIXME add warning logging here
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
		// FIXME ideally we could retry on raw h2 error not grpc error
		// See this ticket for follow up https://github.com/momentohq/client-sdk-go/issues/156
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
