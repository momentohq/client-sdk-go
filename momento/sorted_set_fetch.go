package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetOrder int

const (
	ASCENDING  SortedSetOrder = 0
	DESCENDING SortedSetOrder = 1
)

type SortedSetFetchByIndex struct {
	StartIndex *int32
	EndIndex   *int32
}

type SortedSetFetchByScore struct {
	MinScore *float64
	MaxScore *float64
	Offset   *uint32
	Count    *uint32
}

type SortedSetFetchRequest struct {
	CacheName string
	SetName   string
	Order     SortedSetOrder
	ByIndex   *SortedSetFetchByIndex
	ByScore   *SortedSetFetchByScore

	grpcRequest  *pb.XSortedSetFetchRequest
	grpcResponse *pb.XSortedSetFetchResponse
	response     responses.SortedSetFetchResponse
}

func (r *SortedSetFetchRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchRequest) requestName() string { return "Sorted set fetch" }

func (r *SortedSetFetchRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	if r.ByIndex != nil && r.ByScore != nil {
		return NewMomentoError(
			InvalidArgumentError,
			"Only one of ByIndex or ByScore may be specified",
			nil,
		)
	}

	grpcReq := &pb.XSortedSetFetchRequest{
		SetName:    []byte(r.SetName),
		Order:      pb.XSortedSetFetchRequest_Order(r.Order),
		WithScores: true,
	}

	if r.ByScore != nil {
		by_score := pb.XSortedSetFetchRequest_ByScore{
			ByScore: &pb.XSortedSetFetchRequest_XByScore{
				Min:    &pb.XSortedSetFetchRequest_XByScore_UnboundedMin{},
				Max:    &pb.XSortedSetFetchRequest_XByScore_UnboundedMax{},
				Offset: 0,
				Count:  -1,
			},
		}

		if r.ByScore.MinScore != nil {
			by_score.ByScore.Min = &pb.XSortedSetFetchRequest_XByScore_MinScore{
				MinScore: &pb.XSortedSetFetchRequest_XByScore_XScore{
					Score:     float64(*r.ByScore.MinScore),
					Exclusive: false,
				},
			}
		}

		if r.ByScore.MaxScore != nil {
			by_score.ByScore.Max = &pb.XSortedSetFetchRequest_XByScore_MaxScore{
				MaxScore: &pb.XSortedSetFetchRequest_XByScore_XScore{
					Score:     float64(*r.ByScore.MaxScore),
					Exclusive: false,
				},
			}
		}

		if r.ByScore.Offset != nil {
			by_score.ByScore.Offset = *r.ByScore.Offset
		}

		if r.ByScore.Count != nil {
			by_score.ByScore.Count = int32(*r.ByScore.Count)
		}

		grpcReq.Range = &by_score
	} else {
		// This is the default.
		by_index := pb.XSortedSetFetchRequest_ByIndex{
			ByIndex: &pb.XSortedSetFetchRequest_XByIndex{
				Start: &pb.XSortedSetFetchRequest_XByIndex_UnboundedStart{},
				End:   &pb.XSortedSetFetchRequest_XByIndex_UnboundedEnd{},
			},
		}

		if r.ByIndex.StartIndex != nil {
			by_index.ByIndex.Start = &pb.XSortedSetFetchRequest_XByIndex_InclusiveStartIndex{
				InclusiveStartIndex: *r.ByIndex.StartIndex,
			}
		}

		if r.ByIndex.EndIndex != nil {
			by_index.ByIndex.End = &pb.XSortedSetFetchRequest_XByIndex_ExclusiveEndIndex{
				ExclusiveEndIndex: *r.ByIndex.EndIndex,
			}
		}

		grpcReq.Range = &by_index
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
		r.response = responses.NewSortedSetFetchHit(sortedSetGrpcElementToModel(grpcResp.Found.GetValuesWithScores().Elements))
	case *pb.XSortedSetFetchResponse_Missing:
		r.response = &responses.SortedSetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func sortedSetGrpcElementToModel(grpcSetElements []*pb.XSortedSetElement) []*responses.SortedSetElement {
	var returnList []*responses.SortedSetElement
	for _, element := range grpcSetElements {
		returnList = append(returnList, &responses.SortedSetElement{
			Value: element.Value,
			Score: element.Score,
		})
	}
	return returnList
}
