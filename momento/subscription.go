package momento

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
		rawMsg := new(pb.XTopicItem)
		err := s.grpcClient.RecvMsg(rawMsg)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		f(ctx, &TopicMessageReceiveResponse{
			value: rawMsg.GetValue().GetText(),
		})
	}
}
