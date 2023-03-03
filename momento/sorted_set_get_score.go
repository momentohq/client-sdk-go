package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetGetScoreRequest struct {
	CacheName     string
	SetName       string
	ElementValues []Value

	grpcRequest  *pb.XSortedSetGetScoreRequest
	grpcResponse *pb.XSortedSetGetScoreResponse
	response     responses.SortedSetGetScoreResponse
}

func (r *SortedSetGetScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetGetScoreRequest) requestName() string { return "Sorted set get score" }

func (r *SortedSetGetScoreRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	values, err := momentoValuesToPrimitiveByteList(r.ElementValues)
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

func (r *SortedSetGetScoreRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetGetScore(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetGetScoreRequest) interpretGrpcResponse() error {
	switch grpcResp := r.grpcResponse.SortedSet.(type) {
	case *pb.XSortedSetGetScoreResponse_Found:
		r.response = &responses.SortedSetGetScoreHit{
			Elements: convertSortedSetScoreElement(grpcResp.Found.GetElements()),
		}
	case *pb.XSortedSetGetScoreResponse_Missing:
		r.response = &responses.SortedSetGetScoreMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	return nil
}

func convertSortedSetScoreElement(grpcSetElements []*pb.XSortedSetGetScoreResponse_XSortedSetGetScoreResponsePart) []responses.SortedSetScoreElement {
	var rList []responses.SortedSetScoreElement
	for _, element := range grpcSetElements {
		switch element.Result {
		case pb.ECacheResult_Hit:
			rList = append(rList, &responses.SortedSetScoreHit{Score: element.Score})
		case pb.ECacheResult_Miss:
			rList = append(rList, &responses.SortedSetScoreMiss{})
		default:
			rList = append(rList, &responses.SortedSetScoreInvalid{})
		}
	}
	return rList
}
