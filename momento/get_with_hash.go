package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type GetWithHashRequest struct {
	// Name of the cache to get the item from
	CacheName string
	// string or byte key to be used to store item
	Key Key

	grpcRequest  *pb.XGetWithHashRequest
	grpcResponse *pb.XGetWithHashResponse
	response     responses.GetWithHashResponse
}

func (r *GetWithHashRequest) cacheName() string { return r.CacheName }

func (r *GetWithHashRequest) key() Key { return r.Key }

func (r *GetWithHashRequest) requestName() string { return "GetWithHash" }

func (r *GetWithHashRequest) initGrpcRequest(scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	r.grpcRequest = &pb.XGetWithHashRequest{
		CacheKey: key,
	}

	return nil
}

func (r *GetWithHashRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.GetWithHash(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	r.grpcResponse = resp

	return resp, nil, nil
}

func (r *GetWithHashRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse

	if resp.Result == pb.ECacheResult_Hit {
		r.response = responses.NewGetWithHashHit(resp.CacheBody)
		return nil
	} else if resp.Result == pb.ECacheResult_Miss {
		r.response = &responses.GetWithHashMiss{}
		return nil
	} else {
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
}
