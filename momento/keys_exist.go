package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type KeysExistRequest struct {
	CacheName string
	Keys      []Key

	grpcRequest  *pb.XKeysExistRequest
	grpcResponse *pb.XKeysExistResponse
	response     responses.KeysExistResponse
}

func (r *KeysExistRequest) cacheName() string { return r.CacheName }

func (r *KeysExistRequest) keys() []Key { return r.Keys }

func (r *KeysExistRequest) requestName() string { return "KeysExist" }

func (r *KeysExistRequest) initGrpcRequest(scsDataClient) error {
	var err error
	var keys [][]byte

	if keys, err = prepareKeys(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XKeysExistRequest{
		CacheKeys: keys,
	}

	return nil
}

func (r *KeysExistRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.KeysExist(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	r.grpcResponse = resp

	return resp, nil, nil
}

func (r *KeysExistRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	r.response = responses.NewKeysExistSuccess(resp.Exists)
	return nil
}

func (r *KeysExistRequest) getResponse() interface{} {
	return r.response
}
