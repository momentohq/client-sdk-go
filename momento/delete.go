package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DeleteRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to delete the item.
	Key Key

	grpcRequest  *pb.XDeleteRequest
	grpcResponse *pb.XDeleteResponse
	response     responses.DeleteResponse
}

func (r *DeleteRequest) cacheName() string { return r.CacheName }

func (r *DeleteRequest) key() Key { return r.Key }

func (r *DeleteRequest) requestName() string { return "Delete" }

func (r *DeleteRequest) initGrpcRequest(scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	r.grpcRequest = &pb.XDeleteRequest{CacheKey: key}

	return nil
}

func (r *DeleteRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.Delete(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *DeleteRequest) interpretGrpcResponse() error {
	r.response = &responses.DeleteSuccess{}
	return nil
}
