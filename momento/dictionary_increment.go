package momento

import (
	"context"
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionaryIncrementRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
	Amount         int64
	Ttl            *utils.CollectionTtl

	grpcRequest  *pb.XDictionaryIncrementRequest
	grpcResponse *pb.XDictionaryIncrementResponse
	response     responses.DictionaryIncrementResponse
}

func (r *DictionaryIncrementRequest) cacheName() string { return r.CacheName }

func (r *DictionaryIncrementRequest) field() Value { return r.Field }

func (r *DictionaryIncrementRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *DictionaryIncrementRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *DictionaryIncrementRequest) requestName() string { return "DictionaryFetch" }

func (r *DictionaryIncrementRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var field []byte
	if field, err = prepareField(r); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	if r.Amount == 0 {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Amount must be given and cannot be 0",
			errors.New("invalid argument"),
		)
	}

	r.grpcRequest = &pb.XDictionaryIncrementRequest{
		DictionaryName:  []byte(r.DictionaryName),
		Field:           field,
		Amount:          r.Amount,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}

	return nil
}

func (r *DictionaryIncrementRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionaryIncrement(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionaryIncrementRequest) interpretGrpcResponse() error {
	r.response = responses.NewDictionaryIncrementSuccess(r.grpcResponse.Value)
	return nil
}
