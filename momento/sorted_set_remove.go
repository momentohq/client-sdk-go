package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetRemoveRequest struct {
	CacheName string
	SetName   string
	Values    []Value

	grpcRequest  *pb.XSortedSetRemoveRequest
	grpcResponse *pb.XSortedSetRemoveResponse
	response     responses.SortedSetRemoveResponse
}

func (r *SortedSetRemoveRequest) cacheName() string { return r.CacheName }

func (r *SortedSetRemoveRequest) requestName() string { return "Sorted set remove" }

func (r *SortedSetRemoveRequest) values() []Value { return r.Values }

func (r *SortedSetRemoveRequest) initGrpcRequest(scsDataClient) error {
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

func (r *SortedSetRemoveRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetRemove(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetRemoveRequest) interpretGrpcResponse() error {
	r.response = &responses.SortedSetRemoveSuccess{}
	return nil
}
