package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

func (r *SetLengthRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetLength(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
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
