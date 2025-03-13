package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetGetRankRequest struct {
	CacheName string
	SetName   string
	Value     Value
	Order     SortedSetOrder

	grpcRequest  *pb.XSortedSetGetRankRequest

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

func (r *SortedSetGetRankRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetGetRank(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetGetRankRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XSortedSetGetRankResponse)
	switch rank := myResp.Rank.(type) {
	case *pb.XSortedSetGetRankResponse_ElementRank:
		switch rank.ElementRank.Result {
		case pb.ECacheResult_Hit:
			r.response = responses.SortedSetGetRankHit(rank.ElementRank.Rank)
		case pb.ECacheResult_Miss:
			r.response = &responses.SortedSetGetRankMiss{}
		default:
			return errUnexpectedGrpcResponse(r, myResp)
		}
	case *pb.XSortedSetGetRankResponse_Missing:
		r.response = &responses.SortedSetGetRankMiss{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}
