package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

type ListConcatenateBackRequest struct {
	CacheName           string
	ListName            string
	Values              []Value
	TruncateFrontToSize uint32
	Ttl                 *utils.CollectionTtl

	grpcRequest  *pb.XListConcatenateBackRequest

	response     responses.ListConcatenateBackResponse
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

func (r *ListConcatenateBackRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListConcatenateBack(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListConcatenateBackRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XListConcatenateBackResponse)
	r.response = responses.NewListConcatenateBackSuccess(myResp.ListLength)
	return nil
}
