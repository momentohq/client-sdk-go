package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetRemoveElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []Value
}

func (r *SetRemoveElementsRequest) cacheName() string { return r.CacheName }

func (r *SetRemoveElementsRequest) requestName() string { return "SetRemoveElements" }

func (r *SetRemoveElementsRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	elements, err := momentoValuesToPrimitiveByteList(r.Elements)
	if err != nil {
		return nil, err
	}

	grpcRequest := &pb.XSetDifferenceRequest{
		SetName: []byte(r.SetName),
		Difference: &pb.XSetDifferenceRequest_Subtrahend{
			Subtrahend: &pb.XSetDifferenceRequest_XSubtrahend{
				SubtrahendSet: &pb.XSetDifferenceRequest_XSubtrahend_Set{
					Set: &pb.XSetDifferenceRequest_XSubtrahend_XSet{
						Elements: elements,
					},
				},
			},
		},
	}

	return grpcRequest, nil
}

func (r *SetRemoveElementsRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetDifference(requestMetadata, grpcRequest.(*pb.XSetDifferenceRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetRemoveElementsRequest) interpretGrpcResponse(_ interface{}) (interface{}, error) {
	return &responses.SetRemoveElementsSuccess{}, nil
}

func (c SetRemoveElementsRequest) GetRequestName() string {
	return "SetRemoveElementsRequest"
}
