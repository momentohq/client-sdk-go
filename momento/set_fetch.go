package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetFetchRequest struct {
	CacheName string
	SetName   string

	grpcRequest  *pb.XSetFetchRequest
	grpcResponse *pb.XSetFetchResponse
	response     responses.SetFetchResponse
}

func (r *SetFetchRequest) cacheName() string { return r.CacheName }

func (r *SetFetchRequest) requestName() string { return "SetFetch" }

func (r *SetFetchRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XSetFetchRequest{SetName: []byte(r.SetName)}

	return nil
}

func (r *SetFetchRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetFetch(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetFetchRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.Set.(type) {
	case *pb.XSetFetchResponse_Found:
		r.response = responses.NewSetFetchHit(rtype.Found.Elements)
	case *pb.XSetFetchResponse_Missing:
		r.response = &responses.SetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func (r *SetFetchRequest) getResponse() interface{} { return r.response }
