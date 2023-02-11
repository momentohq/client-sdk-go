package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc/metadata"
)

type SetRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Bytes
	// string ot byte value to be stored.
	Value Bytes
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	TTL time.Duration
}

func (r SetRequest) cacheName() string { return r.CacheName }

func (r SetRequest) key() Bytes { return r.Key }

func (r SetRequest) value() Bytes { return r.Value }

func (r SetRequest) ttl() time.Duration { return r.TTL }

type SetResponse interface {
	isSetResponse()
}

type SetSuccess struct{}

func (SetSuccess) isSetResponse() {}

func (r SetRequest) makeRequest(
	ctx context.Context,
	client DefaultScsClient,
) (SetResponse, error) {
	var err error

	var cache string
	if cache, err = prepareCacheName(r); err != nil {
		return nil, err
	}

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTTL); err != nil {
		return nil, err
	}

	dataClient := client.dataClient

	ctx, cancel := context.WithTimeout(ctx, dataClient.RequestTimeout())
	defer cancel()
	_, err = dataClient.GrpcClient().Set(
		metadata.NewOutgoingContext(ctx, dataClient.CreateNewMetadata(cache)),
		&pb.XSetRequest{
			CacheKey:        key,
			CacheBody:       value,
			TtlMilliseconds: ttl,
		},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return SetSuccess{}, nil
}
