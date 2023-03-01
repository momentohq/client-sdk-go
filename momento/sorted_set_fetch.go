package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

//////// Response

type SortedSetElement struct {
	Value []byte
	Score float64
}
type SortedSetFetchResponse interface {
	isSortedSetFetchResponse()
}

// SortedSetFetchMiss Miss Response to a cache SortedSetFetch api request.
type SortedSetFetchMiss struct{}

func (SortedSetFetchMiss) isSortedSetFetchResponse() {}

// SortedSetFetchHit Hit Response to a cache SortedSetFetch api request.
type SortedSetFetchHit struct {
	Elements []*SortedSetElement
}

func (SortedSetFetchHit) isSortedSetFetchResponse() {}

///// Request

type SortedSetOrder int

const (
	ASCENDING  SortedSetOrder = 0
	DESCENDING SortedSetOrder = 1
)

type SortedSetFetchNumResults interface {
	isSortedSetFetchNumResults()
}

type FetchAllElements struct{}

func (FetchAllElements) isSortedSetFetchNumResults() {}

type FetchLimitedElements struct {
	Limit uint32
}

func (FetchLimitedElements) isSortedSetFetchNumResults() {}

type SortedSetFetchRequest struct {
	CacheName string
	SetName   string
	Order     SortedSetOrder

	grpcRequest  *pb.XSortedSetFetchRequest
	grpcResponse *pb.XSortedSetFetchResponse
	response     SortedSetFetchResponse
}

func (r *SortedSetFetchRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchRequest) requestName() string { return "Sorted set fetch" }

func (r *SortedSetFetchRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	grpcReq := &pb.XSortedSetFetchRequest{
		SetName:    []byte(r.SetName),
		Order:      pb.XSortedSetFetchRequest_Order(r.Order),
		WithScores: true,
		Range: &pb.XSortedSetFetchRequest_ByIndex{
			ByIndex: &pb.XSortedSetFetchRequest_XByIndex{
				Start: &pb.XSortedSetFetchRequest_XByIndex_UnboundedStart{},
				End:   &pb.XSortedSetFetchRequest_XByIndex_UnboundedEnd{},
			},
		},
	}

	r.grpcRequest = grpcReq
	return nil
}

func (r *SortedSetFetchRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetFetch(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetFetchRequest) interpretGrpcResponse() error {
	switch grpcResp := r.grpcResponse.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		r.response = &SortedSetFetchHit{
			Elements: sortedSetGrpcElementToModel(grpcResp.Found.GetValuesWithScores().Elements),
		}
	case *pb.XSortedSetFetchResponse_Missing:
		r.response = &SortedSetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func sortedSetGrpcElementToModel(grpcSetElements []*pb.XSortedSetElement) []*SortedSetElement {
	var returnList []*SortedSetElement
	for _, element := range grpcSetElements {
		returnList = append(returnList, &SortedSetElement{
			Value: element.Value,
			Score: element.Score,
		})
	}
	return returnList
}
