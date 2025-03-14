package momento

import (
	"context"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type GetRequest struct {
	// Name of the cache to get the item from
	CacheName string
	// string or byte key to be used to store item
	Key Key

	grpcRequest *pb.XGetRequest
	response responses.GetResponse
}

func (r *GetRequest) cacheName() string { return r.CacheName }

func (r *GetRequest) key() Key { return r.Key }

func (r *GetRequest) requestName() string { return "Get" }

func (r *GetRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	r.grpcRequest = &pb.XGetRequest{
		CacheKey: key,
	}

	return r.grpcRequest, nil
}

func (r *GetRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.Get(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *GetRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XGetResponse)
	if myResp.Result == pb.ECacheResult_Hit {
		r.response = responses.NewGetHit(myResp.CacheBody)
		return nil
	} else if myResp.Result == pb.ECacheResult_Miss {
		r.response = &responses.GetMiss{}
		return nil
	} else {
		return errUnexpectedGrpcResponse(r, myResp)
	}
}

func (r *GetRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XGetResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
