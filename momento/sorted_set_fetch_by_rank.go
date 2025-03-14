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

	grpcRequest *pb.XSortedSetFetchRequest

	response responses.SortedSetFetchResponse
}

func (r *SortedSetFetchByRankRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchByRankRequest) requestName() string { return "Sorted set fetch" }

func (r *SortedSetFetchByRankRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
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

	r.grpcRequest = grpcReq

	return r.grpcRequest, nil
}

func (r *SortedSetFetchByRankRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetFetch(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetFetchByRankRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XSortedSetFetchResponse)
	switch grpcResp := myResp.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		r.response = responses.NewSortedSetFetchHit(sortedSetByRankGrpcElementToModel(grpcResp.Found.GetValuesWithScores().Elements))
	case *pb.XSortedSetFetchResponse_Missing:
		r.response = &responses.SortedSetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
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

func (r *SortedSetFetchByRankRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XSortedSetFetchResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
