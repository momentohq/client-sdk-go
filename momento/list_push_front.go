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

type ListPushFrontRequest struct {
	CacheName          string
	ListName           string
	Value              Value
	TruncateBackToSize uint32
	Ttl                *utils.CollectionTtl

	grpcRequest *pb.XListPushFrontRequest

	response responses.ListPushFrontResponse
}

func (r *ListPushFrontRequest) cacheName() string { return r.CacheName }

func (r *ListPushFrontRequest) value() Value { return r.Value }

func (r *ListPushFrontRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *ListPushFrontRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *ListPushFrontRequest) requestName() string { return "ListPushFront" }

func (r *ListPushFrontRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	r.grpcRequest = &pb.XListPushFrontRequest{
		ListName:           []byte(r.ListName),
		Value:              value,
		TtlMilliseconds:    ttlMilliseconds,
		RefreshTtl:         refreshTtl,
		TruncateBackToSize: r.TruncateBackToSize,
	}

	return r.grpcRequest, nil
}

func (r *ListPushFrontRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListPushFront(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListPushFrontRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XListPushFrontResponse)
	r.response = responses.NewListPushFrontSuccess(myResp.ListLength)
	return nil
}

func (r *ListPushFrontRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XListPushFrontResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
