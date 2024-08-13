package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetGetScoresRequest struct {
	CacheName string
	SetName   string
	Values    []Value

	grpcRequest  *pb.XSortedSetGetScoreRequest
	grpcResponse *pb.XSortedSetGetScoreResponse
	response     responses.SortedSetGetScoresResponse
}

func (r *SortedSetGetScoresRequest) cacheName() string { return r.CacheName }

func (r *SortedSetGetScoresRequest) requestName() string { return "Sorted set get score" }

func (r *SortedSetGetScoresRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	values, err := momentoValuesToPrimitiveByteList(r.Values)
	if err != nil {
		return err
	}

	resp := &pb.XSortedSetGetScoreRequest{
		SetName: []byte(r.SetName),
		Values:  values,
	}

	r.grpcRequest = resp

	return nil
}

func (r *SortedSetGetScoresRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetGetScore(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetGetScoresRequest) interpretGrpcResponse() error {
	switch grpcResp := r.grpcResponse.SortedSet.(type) {
	case *pb.XSortedSetGetScoreResponse_Found:
		r.response = responses.NewSortedSetGetScoresHit(
			convertSortedSetScoreElement(grpcResp.Found.GetElements()),
			r.grpcRequest.Values,
		)
	case *pb.XSortedSetGetScoreResponse_Missing:
		r.response = &responses.SortedSetGetScoresMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	return nil
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

func (r *SortedSetGetScoresRequest) getResponse() interface{} {
	return r.response
}
