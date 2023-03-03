package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type GetRequest struct {
	// Name of the cache to get the item from
	CacheName string
	// string or byte key to be used to store item
	Key Key

	grpcRequest  *pb.XGetRequest
	grpcResponse *pb.XGetResponse
	response     responses.GetResponse
}

func (r *GetRequest) cacheName() string { return r.CacheName }

func (r *GetRequest) key() Key { return r.Key }

func (r *GetRequest) requestName() string { return "Get" }

func (r *GetRequest) initGrpcRequest(scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	r.grpcRequest = &pb.XGetRequest{
		CacheKey: key,
	}

	return nil
}

func (r *GetRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.Get(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *GetRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse

	if resp.Result == pb.ECacheResult_Hit {
		r.response = responses.NewGetHit(resp.CacheBody)
		return nil
	} else if resp.Result == pb.ECacheResult_Miss {
		r.response = &responses.GetMiss{}
		return nil
	} else {
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
}
