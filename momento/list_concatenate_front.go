package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ListConcatenateFrontRequest struct {
	CacheName          string
	ListName           string
	Values             []Value
	TruncateBackToSize uint32
	Ttl                *utils.CollectionTtl
}

func (r *ListConcatenateFrontRequest) cacheName() string { return r.CacheName }

func (r *ListConcatenateFrontRequest) values() []Value { return r.Values }

func (r *ListConcatenateFrontRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *ListConcatenateFrontRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *ListConcatenateFrontRequest) requestName() string { return "ListConcatenateFront" }

func (r *ListConcatenateFrontRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return nil, err
	}

	var values [][]byte
	if values, err = prepareValues(r); err != nil {
		return nil, err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XListConcatenateFrontRequest{
		ListName:           []byte(r.ListName),
		Values:             values,
		TtlMilliseconds:    ttlMilliseconds,
		RefreshTtl:         refreshTtl,
		TruncateBackToSize: r.TruncateBackToSize,
	}

	return grpcRequest, nil
}

func (r *ListConcatenateFrontRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListConcatenateFront(requestMetadata, grpcRequest.(*pb.XListConcatenateFrontRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListConcatenateFrontRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XListConcatenateFrontResponse)
	return responses.NewListConcatenateFrontSuccess(myResp.ListLength), nil
}
