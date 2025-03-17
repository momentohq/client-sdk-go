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

	response responses.SetLengthResponse
}

func (r *SetLengthRequest) cacheName() string { return r.CacheName }

func (r *SetLengthRequest) requestName() string { return "SetLength" }

func (r *SetLengthRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	if _, err := prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XSetLengthRequest{
		SetName: []byte(r.SetName),
	}

	return grpcRequest, nil
}

func (r *SetLengthRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetLength(requestMetadata, grpcRequest.(*pb.XSetLengthRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetLengthRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XSetLengthResponse)
	switch rtype := myResp.Set.(type) {
	case *pb.XSetLengthResponse_Found:
		r.response = responses.NewSetLengthHit(rtype.Found.Length)
	case *pb.XSetLengthResponse_Missing:
		r.response = &responses.SetLengthMiss{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}

func (r *SetLengthRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XSetLengthResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
