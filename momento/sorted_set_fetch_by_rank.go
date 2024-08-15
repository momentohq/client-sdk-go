package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetFetchByRankRequest struct {
	CacheName string
	SetName   string
	Order     SortedSetOrder
	StartRank *int32
	EndRank   *int32

	grpcRequest  *pb.XSortedSetFetchRequest
	grpcResponse *pb.XSortedSetFetchResponse
	response     responses.SortedSetFetchResponse
}

func (r *SortedSetFetchByRankRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchByRankRequest) requestName() string { return "Sorted set fetch" }

func (r *SortedSetFetchByRankRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
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
			return err
		}
	}

	grpcReq.Range = &byIndex

	r.grpcRequest = grpcReq

	return nil
}

func (r *SortedSetFetchByRankRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetFetch(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetFetchByRankRequest) interpretGrpcResponse() error {
	switch grpcResp := r.grpcResponse.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		r.response = responses.NewSortedSetFetchHit(sortedSetByRankGrpcElementToModel(grpcResp.Found.GetValuesWithScores().Elements))
	case *pb.XSortedSetFetchResponse_Missing:
		r.response = &responses.SortedSetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
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

func (r *SortedSetFetchByRankRequest) getResponse() interface{} {
	return r.response
}
