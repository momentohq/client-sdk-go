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
}

func (r *SortedSetGetRankRequest) cacheName() string { return r.CacheName }

func (r *SortedSetGetRankRequest) requestName() string { return "SortedSetGetRank" }

func (r *SortedSetGetRankRequest) value() Value { return r.Value }

func (r *SortedSetGetRankRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	resp := &pb.XSortedSetGetRankRequest{
		SetName: []byte(r.SetName),
		Value:   value,
		Order:   pb.XSortedSetGetRankRequest_Order(r.Order),
	}
	return resp, nil
}

func (r *SortedSetGetRankRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetGetRank(requestMetadata, grpcRequest.(*pb.XSortedSetGetRankRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetGetRankRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSortedSetGetRankResponse)
	switch rank := myResp.Rank.(type) {
	case *pb.XSortedSetGetRankResponse_ElementRank:
		switch rank.ElementRank.Result {
		case pb.ECacheResult_Hit:
			return responses.SortedSetGetRankHit(rank.ElementRank.Rank), nil
		case pb.ECacheResult_Miss:
			return &responses.SortedSetGetRankMiss{}, nil
		default:
			return nil, errUnexpectedGrpcResponse(r, myResp)
		}
	case *pb.XSortedSetGetRankResponse_Missing:
		return &responses.SortedSetGetRankMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
