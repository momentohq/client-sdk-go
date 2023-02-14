package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// ListPushFrontResponse

type ListPushFrontResponse interface {
	isListPushFrontResponse()
}

type ListPushFrontSuccess struct {
	value uint32
}

func (ListPushFrontSuccess) isListPushFrontResponse() {}

func (resp ListPushFrontSuccess) ListLength() uint32 {
	return resp.value
}

// ListPushFrontRequest

type ListPushFrontRequest struct {
	CacheName          string
	ListName           string
	Value              Bytes
	TruncateBackToSize uint32
	CollectionTTL      utils.CollectionTTL

	grpcRequest  *pb.XListPushFrontRequest
	grpcResponse *pb.XListPushFrontResponse
	response     ListPushFrontResponse
}

func (r ListPushFrontRequest) cacheName() string { return r.CacheName }

func (r ListPushFrontRequest) value() Bytes { return r.Value }

func (r ListPushFrontRequest) ttl() time.Duration { return r.CollectionTTL.Ttl }

func (r ListPushFrontRequest) requestName() string { return "ListPushFront" }

func (r *ListPushFrontRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListPushFrontRequest{
		ListName:           []byte(r.ListName),
		Value:              value,
		TtlMilliseconds:    ttl,
		RefreshTtl:         r.CollectionTTL.RefreshTtl,
		TruncateBackToSize: r.TruncateBackToSize,
	}

	return nil
}

func (r *ListPushFrontRequest) makeGrpcRequest(client scsDataClient, ctx context.Context) (grpcResponse, error) {
	resp, err := client.grpcClient.ListPushFront(ctx, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListPushFrontRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	r.response = ListPushFrontSuccess{value: resp.ListLength}
	return nil
}
