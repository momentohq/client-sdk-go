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

}

func (r *SetContainsElementsRequest) cacheName() string { return r.CacheName }

func (r *SetContainsElementsRequest) requestName() string { return "SetContainsElements" }

func (r *SetContainsElementsRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	var values [][]byte
	for _, v := range r.Elements {
		values = append(values, v.asBytes())
	}

	grpcRequest := &pb.XSetContainsRequest{
		SetName:  []byte(r.SetName),
		Elements: values,
	}

	return grpcRequest, nil
}

func (r *SetContainsElementsRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetContains(requestMetadata, grpcRequest.(*pb.XSetContainsRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetContainsElementsRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSetContainsResponse)
	switch rtype := myResp.Set.(type) {
	case *pb.XSetContainsResponse_Missing:
		return &responses.SetContainsElementsMiss{}, nil
	case *pb.XSetContainsResponse_Found:
		return responses.NewSetContainsElementsHit(rtype.Found.Contains), nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
