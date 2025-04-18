package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetGetScoresRequest struct {
	CacheName string
	SetName   string
	Values    []Value

	grpcResponse *pb.XSortedSetGetScoreResponse
	grpcRequest  *pb.XSortedSetGetScoreRequest
}

func (r *SortedSetGetScoresRequest) cacheName() string { return r.CacheName }

func (r *SortedSetGetScoresRequest) requestName() string { return "SortedSetGetScores" }

func (r *SortedSetGetScoresRequest) initGrpcRequest(scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	values, err := momentoValuesToPrimitiveByteList(r.Values)
	if err != nil {
		return nil, err
	}

	resp := &pb.XSortedSetGetScoreRequest{
		SetName: []byte(r.SetName),
		Values:  values,
	}
	return resp, nil
}

func (r *SortedSetGetScoresRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	// this is used in interpretGrpcResponse
	r.grpcRequest = grpcRequest.(*pb.XSortedSetGetScoreRequest)
	resp, err := client.grpcClient.SortedSetGetScore(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetGetScoresRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	r.grpcResponse = resp.(*pb.XSortedSetGetScoreResponse)
	switch t := r.grpcResponse.SortedSet.(type) {
	case *pb.XSortedSetGetScoreResponse_Found:
		return responses.NewSortedSetGetScoresHit(
			convertSortedSetScoreElement(t.Found.Elements),
			r.grpcRequest.Values,
		), nil
	case *pb.XSortedSetGetScoreResponse_Missing:
		return &responses.SortedSetGetScoresMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
}

func convertSortedSetScoreElement(grpcSetElements []*pb.XSortedSetGetScoreResponse_XSortedSetGetScoreResponsePart) []responses.SortedSetGetScoreResponse {
	var rList []responses.SortedSetGetScoreResponse
	for _, element := range grpcSetElements {
		switch element.Result {
		case pb.ECacheResult_Hit:
			rList = append(rList, responses.NewSortedSetGetScoreHit(element.Score))
		case pb.ECacheResult_Miss:
			rList = append(rList, &responses.SortedSetGetScoreMiss{})
		default:
			rList = append(rList, &responses.SortedSetGetScoreInvalid{})
		}
	}
	return rList
}
