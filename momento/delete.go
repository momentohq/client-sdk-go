package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

func (r *DeleteRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.Delete(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	r.grpcResponse = resp

	return resp, nil, nil
}

func (r *DeleteRequest) interpretGrpcResponse() error {
	r.response = &responses.DeleteSuccess{}
	return nil
}

func (r *DeleteRequest) getResponse() map[string]string { return getMomentoResponseData(r.response) }
