package incubating

import (
	"context"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
	"io"
)

type SubscriptionIFace interface {
	Recv(ctx context.Context, f func(ctx context.Context, m *TopicMessageReceiveResponse)) error
}

type Subscription struct {
	grpcClient grpc.ClientStream
}

func (s *Subscription) Recv(ctx context.Context, f func(ctx context.Context, m *TopicMessageReceiveResponse)) error {
	for {
		rawMsg := new(pb.XSubscriptionItem)
		if err := s.grpcClient.RecvMsg(rawMsg); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		f(ctx, &TopicMessageReceiveResponse{
			// TODO think about user experience for bytes/strings
			value: rawMsg.GetItem().GetValue().GetText(),
		})
	}
}
