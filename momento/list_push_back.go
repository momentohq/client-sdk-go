package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListPushBackRequest struct {
	CacheName           string
	ListName            string
	Value               Value
	TruncateFrontToSize uint32
	Ttl                 *utils.CollectionTtl

	grpcRequest  *pb.XListPushBackRequest

	response     responses.ListPushBackResponse
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

func (r *ListPushBackRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListPushBack(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListPushBackRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XListPushBackResponse)
	r.response = responses.NewListPushBackSuccess(myResp.ListLength)
	return nil
}

func (r *ListPushBackRequest) getResponse() interface{} {
	return r.response
}
