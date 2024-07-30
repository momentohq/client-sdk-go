package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetLengthRequest struct {
	CacheName string
	SetName   string

	grpcRequest  *pb.XSortedSetLengthRequest
	grpcResponse *pb.XSortedSetLengthResponse
	response     responses.SortedSetLengthResponse
}

func (r *SortedSetLengthRequest) cacheName() string { return r.CacheName }

func (r *SortedSetLengthRequest) requestName() string { return "SortedSetLength" }

func (r *SortedSetLengthRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XSortedSetLengthRequest{
		SetName: []byte(r.SetName),
	}

	return nil
}

func (r *SortedSetLengthRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetLength(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SortedSetLengthRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	switch rtype := resp.SortedSet.(type) {
	case *pb.XSortedSetLengthResponse_Found:
		r.response = responses.NewSortedSetLengthHit(rtype.Found.Length)
	case *pb.XSortedSetLengthResponse_Missing:
		r.response = &responses.SortedSetLengthMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func (r *SortedSetLengthRequest) getResponse() map[string]string {
	return getMomentoResponseData(r.response)
}
