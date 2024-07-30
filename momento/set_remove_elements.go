package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetRemoveElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []Value

	grpcRequest  *pb.XSetDifferenceRequest
	grpcResponse *pb.XSetDifferenceResponse
	response     responses.SetRemoveElementsResponse
}

func (r *SetRemoveElementsRequest) cacheName() string { return r.CacheName }

func (r *SetRemoveElementsRequest) requestName() string { return "SetRemoveElements" }

func (r *SetRemoveElementsRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	elements, err := momentoValuesToPrimitiveByteList(r.Elements)
	if err != nil {
		return err
	}

	r.grpcRequest = &pb.XSetDifferenceRequest{
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

	return nil
}

func (r *SetRemoveElementsRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetDifference(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetRemoveElementsRequest) interpretGrpcResponse() error {
	r.response = &responses.SetRemoveElementsSuccess{}
	return nil
}

func (r *SetRemoveElementsRequest) getResponse() map[string]string {
	return getMomentoResponseData(r.response)
}
