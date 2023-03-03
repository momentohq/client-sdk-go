package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
)

type ListConcatenateFrontRequest struct {
	CacheName          string
	ListName           string
	Values             []Value
	TruncateBackToSize uint32
	Ttl                *utils.CollectionTtl

	grpcRequest  *pb.XListConcatenateFrontRequest
	grpcResponse *pb.XListConcatenateFrontResponse
	response     responses.ListConcatenateFrontResponse
}

func (r *ListConcatenateFrontRequest) cacheName() string { return r.CacheName }

func (r *ListConcatenateFrontRequest) values() []Value { return r.Values }

func (r *ListConcatenateFrontRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *ListConcatenateFrontRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

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

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListConcatenateFrontRequest{
		ListName:           []byte(r.ListName),
		Values:             values,
		TtlMilliseconds:    ttlMilliseconds,
		RefreshTtl:         refreshTtl,
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
	r.response = responses.NewListConcatenateFrontSuccess(resp.ListLength)
	return nil
}
