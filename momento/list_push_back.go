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
	Value               Bytes
	TruncateFrontToSize uint32
	CollectionTTL       utils.CollectionTTL

	grpcRequest  *pb.XListPushBackRequest
	grpcResponse *pb.XListPushBackResponse
	response     ListPushBackResponse
}

func (r ListPushBackRequest) cacheName() string { return r.CacheName }

func (r ListPushBackRequest) value() Bytes { return r.Value }

func (r ListPushBackRequest) ttl() time.Duration { return r.CollectionTTL.Ttl }

func (r ListPushBackRequest) requestName() string { return "ListPushBack" }

func (r *ListPushBackRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListPushBackRequest{
		ListName:            []byte(r.ListName),
		Value:               r.Value.AsBytes(),
		TtlMilliseconds:     ttl,
		RefreshTtl:          r.CollectionTTL.RefreshTtl,
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
	r.response = ListPushBackSuccess{value: resp.ListLength}
	return nil
}
