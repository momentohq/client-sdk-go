package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetFetchByScoreRequest struct {
	CacheName string
	SetName   string
	Order     SortedSetOrder
	MinScore  *float64
	MaxScore  *float64
	Offset    *uint32
	Count     *uint32

	grpcRequest  *pb.XSortedSetFetchRequest
	grpcResponse *pb.XSortedSetFetchResponse
	response     responses.SortedSetFetchResponse
}

func (r *SortedSetFetchByScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchByScoreRequest) requestName() string { return "Sorted set fetch" }

func (r *SortedSetFetchByScoreRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	grpcReq := &pb.XSortedSetFetchRequest{
		SetName:    []byte(r.SetName),
		Order:      pb.XSortedSetFetchRequest_Order(r.Order),
		WithScores: true,
	}

	by_score := pb.XSortedSetFetchRequest_ByScore{
		ByScore: &pb.XSortedSetFetchRequest_XByScore{
			Min:    &pb.XSortedSetFetchRequest_XByScore_UnboundedMin{},
			Max:    &pb.XSortedSetFetchRequest_XByScore_UnboundedMax{},
			Offset: 0,
			Count:  -1,
		},
	}

	if r.MinScore != nil {
		by_score.ByScore.Min = &pb.XSortedSetFetchRequest_XByScore_MinScore{
			MinScore: &pb.XSortedSetFetchRequest_XByScore_XScore{
				Score:     float64(*r.MinScore),
				Exclusive: false,
			},
		}
	}

	if r.MaxScore != nil {
		by_score.ByScore.Max = &pb.XSortedSetFetchRequest_XByScore_MaxScore{
			MaxScore: &pb.XSortedSetFetchRequest_XByScore_XScore{
				Score:     float64(*r.MaxScore),
				Exclusive: false,
			},
		}
	}

	if r.Offset != nil {
		by_score.ByScore.Offset = *r.Offset
	}

	if r.Count != nil {
		by_score.ByScore.Count = int32(*r.Count)
	}

	grpcReq.Range = &by_score

	r.grpcRequest = grpcReq

	return nil
}

func (r *SortedSetFetchByScoreRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetFetch(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	r.grpcResponse = resp

	return resp, nil, nil
}

func (r *SortedSetFetchByScoreRequest) interpretGrpcResponse() error {
	switch grpcResp := r.grpcResponse.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		r.response = responses.NewSortedSetFetchHit(sortedSetByScoreGrpcElementToModel(grpcResp.Found.GetValuesWithScores().Elements))
	case *pb.XSortedSetFetchResponse_Missing:
		r.response = &responses.SortedSetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func sortedSetByScoreGrpcElementToModel(grpcSetElements []*pb.XSortedSetElement) []responses.SortedSetBytesElement {
	var returnList []responses.SortedSetBytesElement
	for _, element := range grpcSetElements {
		returnList = append(returnList, responses.SortedSetBytesElement{
			Value: element.Value,
			Score: element.Score,
		})
	}
	return returnList
}

func (r *SortedSetFetchByScoreRequest) getResponse() interface{} {
	return r.response
}
