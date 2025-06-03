package topic_manager_lists

import (
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// Interface for interacting with a pool of TopicGrpcManager objects.
// Implemented by StaticUnaryManagerList.
type TopicManagerList interface {
	GetNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr)
	Close()
}

// StaticUnaryManagerList manages a static pool of gRPC channels for unary pubsub requests.
type StaticUnaryManagerList struct {
	grpcManagers []*grpcmanagers.TopicGrpcManager
	managerIndex atomic.Uint64
	logger       logger.MomentoLogger
}

// Each grpc connection can multiplex 100 subscribe/publish requests.
// Publish requests will queue up on client while waiting for in-flight requests to complete if
// the number of concurrent requests exceeds numUnaryChannels*100, but will eventually complete.
// Therefore we can just round-robin the unaryTopicManagers, no need to keep track of how many
// publish requests are in flight on each one.
func (list *StaticUnaryManagerList) GetNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	nextManagerIndex := list.managerIndex.Add(1)
	return list.grpcManagers[nextManagerIndex%uint64(len(list.grpcManagers))], nil
}

func (list *StaticUnaryManagerList) Close() {
	for _, topicManager := range list.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			list.logger.Error("Error closing topic manager: %v", err)
		}
	}
}

func NewStaticUnaryManagerList(request *models.TopicStreamGrpcManagerRequest, numUnaryChannels uint32, logger logger.MomentoLogger) (*StaticUnaryManagerList, momentoerrors.MomentoSvcErr) {
	unaryTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	for i := 0; uint32(i) < numUnaryChannels; i++ {
		unaryTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(request)
		if err != nil {
			return nil, err
		}
		unaryTopicManagers = append(unaryTopicManagers, unaryTopicManager)
	}
	return &StaticUnaryManagerList{
		grpcManagers: unaryTopicManagers,
		logger:       logger,
	}, nil
}
