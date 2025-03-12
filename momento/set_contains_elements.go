package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetContainsElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []Value

	grpcRequest  *pb.XSetContainsRequest
	grpcResponse *pb.XSetContainsResponse
	response     responses.SetContainsElementsResponse
}

func (r *SetContainsElementsRequest) cacheName() string { return r.CacheName }

func (r *SetContainsElementsRequest) requestName() string { return "SetContainsElements" }

func (r *SetContainsElementsRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var values [][]byte
	for _, v := range r.Elements {
		values = append(values, v.asBytes())
	}

	r.grpcRequest = &pb.XSetContainsRequest{
		SetName:  []byte(r.SetName),
		Elements: values,
	}

	return nil
}

func (r *SetContainsElementsRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetContains(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *SetContainsElementsRequest) interpretGrpcResponse(_ interface{}) error {
	switch rtype := r.grpcResponse.Set.(type) {
	case *pb.XSetContainsResponse_Missing:
		r.response = &responses.SetContainsElementsMiss{}
	case *pb.XSetContainsResponse_Found:
		r.response = responses.NewSetContainsElementsHit(rtype.Found.Contains)
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func (r *SetContainsElementsRequest) getResponse() interface{} {
	return r.response
}
