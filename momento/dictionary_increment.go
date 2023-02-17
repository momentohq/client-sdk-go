package momento

import (
	"context"
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionaryIncrementResponse

type DictionaryIncrementResponse interface {
	isDictionaryIncrementResponse()
}

type DictionaryIncrementSuccess struct {
	value int64
}

func (DictionaryIncrementSuccess) isDictionaryIncrementResponse() {}

func (resp DictionaryIncrementSuccess) ValueUint64() int64 {
	return resp.value
}

// DictionaryIncrementRequest

type DictionaryIncrementRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
	Amount         int64
	CollectionTTL  utils.CollectionTTL

	grpcRequest  *pb.XDictionaryIncrementRequest
	grpcResponse *pb.XDictionaryIncrementResponse
	response     DictionaryIncrementResponse
}

func (r *DictionaryIncrementRequest) cacheName() string { return r.CacheName }

func (r *DictionaryIncrementRequest) field() Value { return r.Field }

func (r *DictionaryIncrementRequest) ttl() time.Duration { return r.CollectionTTL.Ttl }

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

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
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
		TtlMilliseconds: ttl,
		RefreshTtl:      r.CollectionTTL.RefreshTtl,
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
	r.response = &DictionaryIncrementSuccess{value: r.grpcResponse.Value}
	return nil
}
