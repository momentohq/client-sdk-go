package momento

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	client_sdk_go "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc/metadata"
)

type GetRequest struct {
	// Name of the cache to get the item from
	CacheName string
	// string or byte key to be used to store item
	Key Bytes
}

func (r GetRequest) cacheName() string { return r.CacheName }

func (r GetRequest) key() Bytes { return r.Key }

type GetResponse interface {
	isGetResponse()
}

// Miss response to a cache Get api request.
type GetMiss struct{}

func (GetMiss) isGetResponse() {}

// Hit response to a cache Get api request.
type GetHit struct {
	value []byte
}

func (GetHit) isGetResponse() {}

// ValueString Returns value stored in cache as string if there was Hit. Returns an empty string otherwise.
func (resp GetHit) ValueString() string {
	return string(resp.value)
}

// ValueByte Returns value stored in cache as bytes if there was Hit. Returns nil otherwise.
func (resp GetHit) ValueByte() []byte {
	return resp.value
}

func (r GetRequest) makeRequest(ctx context.Context, client DefaultScsClient) (GetResponse, error) {
	var err error

	var cache string
	if cache, err = prepareCacheName(r); err != nil {
		return nil, err
	}

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	dataClient := client.dataClient

	ctx, cancel := context.WithTimeout(ctx, dataClient.RequestTimeout())
	defer cancel()

	resp, err := dataClient.GrpcClient().Get(
		metadata.NewOutgoingContext(ctx, dataClient.CreateNewMetadata(cache)),
		&client_sdk_go.XGetRequest{
			CacheKey: key,
		},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	if resp.Result == client_sdk_go.ECacheResult_Hit {
		return &GetHit{value: resp.CacheBody}, nil
	} else if resp.Result == client_sdk_go.ECacheResult_Miss {
		return &GetMiss{}, nil
	} else {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InternalServerError,
			fmt.Sprintf(
				"CacheService returned an unexpected result: %v for operation: %s with message: %s",
				resp.Result, "GET", resp.Message,
			),
			nil,
		)
	}
}
