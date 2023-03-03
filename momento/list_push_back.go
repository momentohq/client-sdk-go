package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// ListPushBackResponse

type ListPushBackResponse interface {
	isListPushBackResponse()
}

type ListPushBackSuccess struct {
	value uint32
}

func (ListPushBackSuccess) isListPushBackResponse() {}

func (resp ListPushBackSuccess) ListLength() uint32 {
	return resp.value
}

// ListPushBackRequest

type ListPushBackRequest struct {
	CacheName           string
	ListName            string
	Value               Value
	TruncateFrontToSize uint32
	Ttl                 *utils.CollectionTtl

	grpcRequest  *pb.XListPushBackRequest
	grpcResponse *pb.XListPushBackResponse
	response     ListPushBackResponse
}

func (r *ListPushBackRequest) cacheName() string { return r.CacheName }

func (r *ListPushBackRequest) value() Value { return r.Value }

func (r *ListPushBackRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *ListPushBackRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *ListPushBackRequest) requestName() string { return "ListPushBack" }

func (r *ListPushBackRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListPushBackRequest{
		ListName:            []byte(r.ListName),
		Value:               value,
		TtlMilliseconds:     ttlMilliseconds,
		RefreshTtl:          refreshTtl,
		TruncateFrontToSize: r.TruncateFrontToSize,
	}

	return nil
}

func (r *ListPushBackRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ListPushBack(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListPushBackRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	r.response = &ListPushBackSuccess{value: resp.ListLength}
	return nil
}
