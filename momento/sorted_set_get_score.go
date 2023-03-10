package momento

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type SortedSetGetScoreRequest struct {
	CacheName string
	SetName   string
	Value     Value

	grpcRequest  *pb.XSortedSetGetScoreRequest
	grpcResponse *pb.XSortedSetGetScoreResponse
	response     responses.SortedSetGetScoresResponse
}
