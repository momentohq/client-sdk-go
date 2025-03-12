package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfPresentAndNotEqualRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// string or byte value to compare with the existing value in the cache.
	NotEqual Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration

	grpcRequest  *pb.XSetIfRequest

	response     responses.SetIfPresentAndNotEqualResponse
}

func (r *SetIfPresentAndNotEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfPresentAndNotEqualRequest) key() Key { return r.Key }

func (r *SetIfPresentAndNotEqualRequest) value() Value { return r.Value }

func (r *SetIfPresentAndNotEqualRequest) notEqual() Value { return r.NotEqual }

func (r *SetIfPresentAndNotEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfPresentAndNotEqualRequest) requestName() string { return "SetIfNotExists" }

func (r *SetIfPresentAndNotEqualRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var notEqual []byte
	if notEqual, err = prepareNotEqual(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return err
	}

	var condition = &pb.XSetIfRequest_PresentAndNotEqual{
		PresentAndNotEqual: &pb.PresentAndNotEqual{
			ValueToCheck: notEqual,
		},
	}
	r.grpcRequest = &pb.XSetIfRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       condition,
	}

	return nil
}

func (r *SetIfPresentAndNotEqualRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetIfPresentAndNotEqualRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XSetIfResponse)
	var theResponse responses.SetIfPresentAndNotEqualResponse

	switch myResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		theResponse = &responses.SetIfPresentAndNotEqualStored{}
	case *pb.XSetIfResponse_NotStored:
		theResponse = &responses.SetIfPresentAndNotEqualNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}

	r.response = theResponse
	return nil
}

func (r *SetIfPresentAndNotEqualRequest) getResponse() interface{} {
	return r.response
}
