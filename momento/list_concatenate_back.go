package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

// ListConcatenateFrontResponse

type ListConcatenateBackResponse interface {
	isListConcatenateBackResponse()
}

type ListConcatenateBackSuccess struct {
	listLength uint32
}

func (ListConcatenateBackSuccess) isListConcatenateBackResponse() {}

func (resp ListConcatenateBackSuccess) ListLength() uint32 {
	return resp.listLength
}

// ListConcatenateBackRequest

type ListConcatenateBackRequest struct {
	CacheName           string
	ListName            string
	Values              []Value
	TruncateFrontToSize uint32
	Ttl                 *utils.CollectionTtl

	grpcRequest  *pb.XListConcatenateBackRequest
	grpcResponse *pb.XListConcatenateBackResponse
	response     ListConcatenateBackResponse
}

func (r *ListConcatenateBackRequest) cacheName() string { return r.CacheName }

func (r *ListConcatenateBackRequest) values() []Value { return r.Values }

func (r *ListConcatenateBackRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *ListConcatenateBackRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *ListConcatenateBackRequest) requestName() string { return "ListConcatenateBack" }

func (r *ListConcatenateBackRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	var values [][]byte
	if values, err = prepareValues(r); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListConcatenateBackRequest{
		ListName:            []byte(r.ListName),
		Values:              values,
		TtlMilliseconds:     ttlMilliseconds,
		RefreshTtl:          refreshTtl,
		TruncateFrontToSize: r.TruncateFrontToSize,
	}

	return nil
}

func (r *ListConcatenateBackRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ListConcatenateBack(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListConcatenateBackRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	r.response = &ListConcatenateBackSuccess{listLength: resp.ListLength}
	return nil
}
