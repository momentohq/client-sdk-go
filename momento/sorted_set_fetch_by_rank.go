package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetFetchByRankRequest struct {
	CacheName string
	SetName   string
	Order     SortedSetOrder
	StartRank *int32
	EndRank   *int32

}

func (r *SortedSetFetchByRankRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchByRankRequest) requestName() string { return "SortedSetFetchByRank" }

func (r *SortedSetFetchByRankRequest) initGrpcRequest(scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	grpcReq := &pb.XSortedSetFetchRequest{
		SetName:    []byte(r.SetName),
		Order:      pb.XSortedSetFetchRequest_Order(r.Order),
		WithScores: true,
	}

	// This is the default: fetch everything in ascending order.
	byIndex := pb.XSortedSetFetchRequest_ByIndex{
		ByIndex: &pb.XSortedSetFetchRequest_XByIndex{
			Start: &pb.XSortedSetFetchRequest_XByIndex_UnboundedStart{},
			End:   &pb.XSortedSetFetchRequest_XByIndex_UnboundedEnd{},
		},
	}

	startForValidation := int32(0)
	if r.StartRank != nil {

		byIndex.ByIndex.Start = &pb.XSortedSetFetchRequest_XByIndex_InclusiveStartIndex{
			InclusiveStartIndex: *r.StartRank,
		}
		startForValidation = *r.StartRank
	}

	if r.EndRank != nil {
		byIndex.ByIndex.End = &pb.XSortedSetFetchRequest_XByIndex_ExclusiveEndIndex{
			ExclusiveEndIndex: *r.EndRank,
		}
		endForValidation := *r.EndRank
		if err := validateSortedSetRanks(startForValidation, endForValidation); err != nil {
			return nil, err
		}
	}

	grpcReq.Range = &byIndex
	return grpcReq, nil
}

func (r *SortedSetFetchByRankRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetFetch(requestMetadata, grpcRequest.(*pb.XSortedSetFetchRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetFetchByRankRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSortedSetFetchResponse)
	switch grpcResp := myResp.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		return responses.NewSortedSetFetchHit(sortedSetByRankGrpcElementToModel(grpcResp.Found.GetValuesWithScores().Elements)), nil
	case *pb.XSortedSetFetchResponse_Missing:
		return &responses.SortedSetFetchMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}

func sortedSetByRankGrpcElementToModel(grpcSetElements []*pb.XSortedSetElement) []responses.SortedSetBytesElement {
	var returnList []responses.SortedSetBytesElement
	for _, element := range grpcSetElements {
		returnList = append(returnList, responses.SortedSetBytesElement{
			Value: element.Value,
			Score: element.Score,
		})
	}
	return returnList
}
