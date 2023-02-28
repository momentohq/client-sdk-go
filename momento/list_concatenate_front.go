package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

// ListConcatenateFrontResponse

type ListConcatenateFrontResponse interface {
	isListConcatenateFrontResponse()
}

type ListConcatenateFrontSuccess struct {
	listLength uint32
}

func (ListConcatenateFrontSuccess) isListConcatenateFrontResponse() {}

func (resp ListConcatenateFrontSuccess) ListLength() uint32 {
	return resp.listLength
}

// ListConcatenateFrontRequest

type ListConcatenateFrontRequest struct {
	CacheName          string
	ListName           string
	Values             []Value
	TruncateBackToSize uint32
	Ttl                utils.CollectionTtl

	grpcRequest  *pb.XListConcatenateFrontRequest
	grpcResponse *pb.XListConcatenateFrontResponse
	response     ListConcatenateFrontResponse
}

func (r *ListConcatenateFrontRequest) cacheName() string { return r.CacheName }

func (r *ListConcatenateFrontRequest) values() []Value { return r.Values }

func (r *ListConcatenateFrontRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *ListConcatenateFrontRequest) requestName() string { return "ListConcatenateFront" }

func (r *ListConcatenateFrontRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	var values [][]byte
	if values, err = prepareValues(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListConcatenateFrontRequest{
		ListName:           []byte(r.ListName),
		Values:             values,
		TtlMilliseconds:    ttl,
		RefreshTtl:         r.Ttl.RefreshTtl,
		TruncateBackToSize: r.TruncateBackToSize,
	}

	return nil
}

func (r *ListConcatenateFrontRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ListConcatenateFront(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListConcatenateFrontRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	r.response = &ListConcatenateFrontSuccess{listLength: resp.ListLength}
	return nil
}
