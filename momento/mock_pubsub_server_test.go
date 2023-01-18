package momento

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type TestPubSubServer struct {
	pb.UnimplementedPubsubServer
}

func newMockPubSubServer() {
	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPubsubServer(s, &TestPubSubServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func (TestPubSubServer) Subscribe(req *pb.XSubscriptionRequest, server pb.Pubsub_SubscribeServer) error {
	count := 0
	for {
		err := server.SendMsg(&pb.XTopicItem{
			TopicSequenceNumber: uint64(count),
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Text{
					Text: fmt.Sprintf("count %d", count),
				},
			},
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
		count += 1
	}
}
