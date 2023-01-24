package incubating

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type TestMomentoLocalServer struct {
	pb.UnimplementedPubsubServer
	// TODO make this more sophisticated to support multiple subscriptions right now just support one global channel to start
	basicMessageChannel chan string
}

func newMomentoLocalTestServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPubsubServer(s, &TestMomentoLocalServer{
		basicMessageChannel: make(chan string),
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (t TestMomentoLocalServer) Publish(ctx context.Context, req *pb.XPublishRequest) (*pb.XEmpty, error) {
	t.basicMessageChannel <- req.Value.String() // TODO think about bytes vs strings
	return &pb.XEmpty{}, nil
}
func (t TestMomentoLocalServer) Subscribe(req *pb.XSubscriptionRequest, server pb.Pubsub_SubscribeServer) error {
	count := 0
	for msg := range t.basicMessageChannel {
		err := server.SendMsg(&pb.XSubscriptionItem{
			Kind: &pb.XSubscriptionItem_Item{
				Item: &pb.XTopicItem{
					TopicSequenceNumber: uint64(count),
					Value: &pb.XTopicValue{
						Kind: &pb.XTopicValue_Text{
							Text: msg,
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		count += 1
	}
	return nil
}
