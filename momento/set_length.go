package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetLengthRequest struct {
	CacheName string
	SetName   string

	grpcRequest  *pb.XSetLengthRequest
	grpcResponse *pb.XSetLengthResponse
	response     responses.SetLengthResponse
}

func (r *SetLengthRequest) cacheName() string { return r.CacheName }

func (r *SetLengthRequest) requestName() string { return "SetLength" }

func (r *SetLengthRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XSetLengthRequest{
		SetName: []byte(r.SetName),
	}

	return nil
}

func (r *SetLengthRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetLength(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetLengthRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	switch rtype := resp.Set.(type) {
	case *pb.XSetLengthResponse_Found:
		r.response = responses.NewSetLengthHit(rtype.Found.Length)
	case *pb.XSetLengthResponse_Missing:
		r.response = &responses.SetLengthMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
