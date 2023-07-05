package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetLengthByScoreRequest struct {
	CacheName string
	SetName   string
	MinScore  *float64
	MaxScore  *float64

	grpcRequest  *pb.XSortedSetLengthByScoreRequest
	grpcResponse *pb.XSortedSetLengthByScoreResponse
	response     responses.SortedSetLengthByScoreResponse
}

func (r *SortedSetLengthByScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetLengthByScoreRequest) requestName() string { return "SortedSetLengthByScore" }

func (r *SortedSetLengthByScoreRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XSortedSetLengthByScoreRequest{
		SetName: []byte(r.SetName),
	}

	return nil
}

func (r *SortedSetLengthByScoreRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetLengthByScore(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SortedSetLengthByScoreRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	switch rtype := resp.SortedSet.(type) {
	case *pb.XSortedSetLengthByScoreResponse_Found:
		r.response = responses.NewSortedSetLengthByScoreHit(rtype.Found.Length)
	case *pb.XSortedSetLengthByScoreResponse_Missing:
		r.response = &responses.SortedSetLengthByScoreMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
