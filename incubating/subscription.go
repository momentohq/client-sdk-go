package incubating

import (
	"context"
	"io"

	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type SubscriptionIFace interface {
	Recv(ctx context.Context, f func(ctx context.Context, m TopicValue)) error
}

type Subscription struct {
	grpcClient grpc.ClientStream
}

func (s *Subscription) Recv(ctx context.Context, f func(ctx context.Context, m TopicValue)) error {
	for {
		rawMsg := new(pb.XSubscriptionItem)
		if err := s.grpcClient.RecvMsg(rawMsg); err != nil {
			if err == io.EOF {
				// TODO think about retry and re-establish more
				return nil
			}
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		switch typedMsg := rawMsg.Kind.(type) {
		case *pb.XSubscriptionItem_Discontinuity:
			// Don't pass discontinuity messages back to user for now
			// TODO decide how want to notify client
		case *pb.XSubscriptionItem_Item:
			switch subscriptionItem := typedMsg.Item.Value.Kind.(type) {
			case *pb.XTopicValue_Text:
				f(ctx, &TopicValueString{
					Text: subscriptionItem.Text,
				})
			case *pb.XTopicValue_Binary:
				f(ctx, &TopicValueBytes{
					Bytes: subscriptionItem.Binary,
				})
			}
		}
	}
}
