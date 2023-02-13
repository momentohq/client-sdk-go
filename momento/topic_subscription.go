package momento

import (
	"io"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
)

type TopicSubscription interface {
	Item() (TopicValue, error)
}

type topicSubscription struct {
	grpcClient grpc.ClientStream
}

func (s topicSubscription) Item() (TopicValue, error) {
	retryCount := 0
	maxRetries := 2
	for {
		rawMsg := new(pb.XSubscriptionItem)
		if err := s.grpcClient.RecvMsg(rawMsg); err != nil {
			if err == io.EOF {
				// TODO think about retry and re-establish more
				return nil, nil
			}
			return nil, err
		}
		switch typedMsg := rawMsg.Kind.(type) {
		case *pb.XSubscriptionItem_Discontinuity:
			retryCount++
			if retryCount > maxRetries {
				return nil, nil
			}
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
		}
	}
}
