package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	client_sdk_go "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc/metadata"
)

type DeleteRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to delete the item.
	Key Bytes
}

func (r DeleteRequest) cacheName() string { return r.CacheName }

func (r DeleteRequest) key() Bytes { return r.Key }

type DeleteResponse interface {
	isDeleteResponse()
}

type DeleteSuccess struct{}

func (DeleteSuccess) isDeleteResponse() {}

func (r DeleteRequest) makeRequest(
	ctx context.Context,
	client scsDataClient,
) (DeleteResponse, error) {
	var err error

	var cache string
	if cache, err = prepareCacheName(r); err != nil {
		return nil, err
	}

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	_, err = client.grpcClient.Delete(
		metadata.NewOutgoingContext(ctx, client.CreateNewMetadata(cache)),
		&client_sdk_go.XDeleteRequest{CacheKey: key},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	return DeleteSuccess{}, nil
}
