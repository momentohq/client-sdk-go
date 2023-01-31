package incubating

import (
	"context"
	"io"

	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type SubscriptionIFace interface {
	Recv(ctx context.Context, f func(ctx context.Context, m TopicMessage)) error
}

type Subscription struct {
	grpcClient grpc.ClientStream
}

func (s *Subscription) Recv(ctx context.Context, f func(ctx context.Context, m TopicMessage)) error {
	for {
		rawMsg := new(pb.XSubscriptionItem)
		if err := s.grpcClient.RecvMsg(rawMsg); err != nil {
			if err == io.EOF {
				// TODO think about retry and re-establish more
				return nil
			}
			return err
		}

		// Don't pass discontinuity messages back to user for now
		if rawMsg.GetItem() != nil {
			if rawMsg.GetItem().GetValue().GetBinary() != nil {
				f(ctx, &TopicMessageBytes{
					Value: rawMsg.GetItem().GetValue().GetBinary(),
				})
			} else if rawMsg.GetItem().GetValue().GetText() != "" {
				f(ctx, &TopicMessageString{
					Value: rawMsg.GetItem().GetValue().GetText(),
				})
			}
		}
	}
}
