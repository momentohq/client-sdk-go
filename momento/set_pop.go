package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

func (r *SetPopRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetPop(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
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

func (r *SetPopRequest) getResponse() interface{} { return r.response }
