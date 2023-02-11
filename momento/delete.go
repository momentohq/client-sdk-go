package momento

import (
	"context"

	client_sdk_go "github.com/momentohq/client-sdk-go/internal/protos"
)

////// DeleteResponse //////

type DeleteResponse interface {
	isDeleteResponse()
}

type DeleteSuccess struct{}

func (DeleteSuccess) isDeleteResponse() {}

////// DeleteRequest //////

type DeleteRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to delete the item.
	Key Bytes

	grpcRequest  *client_sdk_go.XDeleteRequest
	grpcResponse *client_sdk_go.XDeleteResponse
	response     DeleteResponse
}

func (r DeleteRequest) cacheName() string { return r.CacheName }

func (r DeleteRequest) key() Bytes { return r.Key }

func (r DeleteRequest) requestName() string { return "Delete" }

func (r *DeleteRequest) initGrpcRequest(scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	r.grpcRequest = &client_sdk_go.XDeleteRequest{CacheKey: key}

	return nil
}

func (r *DeleteRequest) makeGrpcRequest(client scsDataClient, metadata context.Context) (grpcResponse, error) {
	resp, err := client.grpcClient.Delete(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *DeleteRequest) interpretGrpcResponse() error {
	r.response = DeleteSuccess{}
	return nil
}
