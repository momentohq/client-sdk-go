package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetGetRankRequest struct {
	CacheName string
	SetName   string
	Value     Value
	Order     SortedSetOrder

	grpcRequest  *pb.XSortedSetGetRankRequest
	grpcResponse *pb.XSortedSetGetRankResponse
	response     responses.SortedSetGetRankResponse
}

func (r *SortedSetGetRankRequest) cacheName() string { return r.CacheName }

func (r *SortedSetGetRankRequest) requestName() string { return "Sorted set get rank" }

func (r *SortedSetGetRankRequest) value() Value { return r.Value }

func (r *SortedSetGetRankRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	resp := &pb.XSortedSetGetRankRequest{
		SetName: []byte(r.SetName),
		Value:   value,
		Order:   pb.XSortedSetGetRankRequest_Order(r.Order),
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

	var resp responses.SortedSetGetRankResponse
	switch rank := grpcResp.Rank.(type) {
	case *pb.XSortedSetGetRankResponse_ElementRank:
		switch rank.ElementRank.Result {
		case pb.ECacheResult_Hit:
			resp = responses.SortedSetGetRankHit(rank.ElementRank.Rank)
		case pb.ECacheResult_Miss:
			resp = &responses.SortedSetGetRankMiss{}
		default:
			return errUnexpectedGrpcResponse(r, r.grpcResponse)
		}
	case *pb.XSortedSetGetRankResponse_Missing:
		resp = &responses.SortedSetGetRankMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp

	return nil
}
