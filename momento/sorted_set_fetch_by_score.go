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
}

func (r *SortedSetFetchByScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetFetchByScoreRequest) requestName() string { return "SortedSetFetchByScore" }

func (r *SortedSetFetchByScoreRequest) initGrpcRequest(scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
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
	return grpcReq, nil
}

func (r *SortedSetFetchByScoreRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetFetch(requestMetadata, grpcRequest.(*pb.XSortedSetFetchRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetFetchByScoreRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSortedSetFetchResponse)
	switch grpcResp := myResp.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		return responses.NewSortedSetFetchHit(sortedSetByScoreGrpcElementToModel(grpcResp.Found.GetValuesWithScores().Elements)), nil
	case *pb.XSortedSetFetchResponse_Missing:
		return &responses.SortedSetFetchMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
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
