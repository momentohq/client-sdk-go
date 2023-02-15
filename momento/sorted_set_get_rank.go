package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

///////// Response

type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMiss Miss Response to a cache SortedSetGetRank api request.
type SortedSetGetRankMiss struct{}

func (SortedSetGetRankMiss) isSortedSetGetRankResponse() {}

// SortedSetGetRankHit Hit Response to a cache SortedSetGetRank api request.
type SortedSetGetRankHit struct {
	Rank uint64
}

func (SortedSetGetRankHit) isSortedSetGetRankResponse() {}

///////// Request

type SortedSetGetRankRequest struct {
	CacheName   string
	SetName     string
	ElementName Bytes

	grpcRequest  *pb.XSortedSetGetRankRequest
	grpcResponse *pb.XSortedSetGetRankResponse
	response     SortedSetGetRankResponse
}

func (r *SortedSetGetRankRequest) cacheName() string { return r.CacheName }

func (r *SortedSetGetRankRequest) requestName() string { return "Sorted set get rank" }

func (r *SortedSetGetRankRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	resp := &pb.XSortedSetGetRankRequest{
		SetName:     []byte(r.SetName),
		ElementName: r.ElementName.asBytes(),
	}

	r.grpcRequest = resp

	return nil
}

func (r *SortedSetGetRankRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetGetRank(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetGetRankRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse

	var resp SortedSetGetRankResponse
	switch rank := grpcResp.Rank.(type) {
	case *pb.XSortedSetGetRankResponse_ElementRank:
		resp = &SortedSetGetRankHit{
			Rank: rank.ElementRank.Rank,
		}
	case *pb.XSortedSetGetRankResponse_Missing:
		resp = &SortedSetGetRankMiss{}
	default:
		return errUnexpectedGrpcResponse
	}

	r.response = resp

	return nil
}
