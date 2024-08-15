package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type ItemGetTtlRequest struct {
	CacheName string
	Key       Key

	grpcRequest  *pb.XItemGetTtlRequest
	grpcResponse *pb.XItemGetTtlResponse
	response     responses.ItemGetTtlResponse
}

func (r *ItemGetTtlRequest) cacheName() string { return r.CacheName }

func (r *ItemGetTtlRequest) key() Key { return r.Key }

func (r *ItemGetTtlRequest) requestName() string { return "ItemGetTypeTL" }

func (r *ItemGetTtlRequest) initGrpcRequest(scsDataClient) error {
	var err error
	var key []byte

	if key, err = prepareKey(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XItemGetTtlRequest{CacheKey: key}

	return nil
}

func (r *ItemGetTtlRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ItemGetTtl(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *ItemGetTtlRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse

	switch grpcResp.Result.(type) {
	case *pb.XItemGetTtlResponse_Found:
		r.response = responses.NewItemGetTtlHit(grpcResp.GetFound().GetRemainingTtlMillis())
		return nil
	case *pb.XItemGetTtlResponse_Missing:
		r.response = &responses.ItemGetTtlMiss{}
		return nil
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
}

func (r *ItemGetTtlRequest) getResponse() interface{} {
	return r.response
}
