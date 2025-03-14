package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetRemoveElementsRequest struct {
	CacheName string
	SetName   string
	Values    []Value

	grpcRequest  *pb.XSortedSetRemoveRequest
	grpcResponse *pb.XSortedSetRemoveResponse
	response     responses.SortedSetRemoveElementsResponse
}

func (r *SortedSetRemoveElementsRequest) cacheName() string { return r.CacheName }

func (r *SortedSetRemoveElementsRequest) requestName() string { return "SortedSetRemoveElements" }

func (r *SortedSetRemoveElementsRequest) values() []Value { return r.Values }

func (r *SortedSetRemoveElementsRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var valuesToRemove [][]byte
	if valuesToRemove, err = prepareValues(r); err != nil {
		return err
	}

	grpcReq := &pb.XSortedSetRemoveRequest{
		SetName: []byte(r.SetName),
	}

	grpcReq.RemoveElements = &pb.XSortedSetRemoveRequest_Some{
		Some: &pb.XSortedSetRemoveRequest_XSome{Values: valuesToRemove},
	}

	r.grpcRequest = grpcReq

	return nil
}

func (r *SortedSetRemoveElementsRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetRemove(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	r.grpcResponse = resp

	return resp, nil, nil
}

func (r *SortedSetRemoveElementsRequest) interpretGrpcResponse() error {
	r.response = &responses.SortedSetRemoveElementsSuccess{}
	return nil
}
