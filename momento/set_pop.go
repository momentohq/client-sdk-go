package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetPopRequest struct {
	CacheName string
	SetName   string
	Count     *uint32

	grpcRequest  *pb.XSetPopRequest
	grpcResponse *pb.XSetPopResponse
	response     responses.SetPopResponse
}

func (r *SetPopRequest) cacheName() string { return r.CacheName }

func (r *SetPopRequest) requestName() string { return "SetPop" }

func (r *SetPopRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var count uint32 = 1
	if r.Count != nil {
		count = uint32(*r.Count)
	}

	r.grpcRequest = &pb.XSetPopRequest{
		SetName: []byte(r.SetName),
		Count:   count,
	}

	return nil
}

func (r *SetPopRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetPop(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetPopRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.Set.(type) {
	case *pb.XSetPopResponse_Found:
		r.response = responses.NewSetPopHit(rtype.Found.Elements)
	case *pb.XSetPopResponse_Missing:
		r.response = &responses.SetPopMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
