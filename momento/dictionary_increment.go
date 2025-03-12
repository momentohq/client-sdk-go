package momento

import (
	"context"
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

func (r *DictionaryIncrementRequest) requestName() string { return "DictionaryIncrement" }

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

func (r *DictionaryIncrementRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionaryIncrement(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *DictionaryIncrementRequest) interpretGrpcResponse(_ interface{}) error {
	r.response = responses.NewDictionaryIncrementSuccess(r.grpcResponse.Value)
	return nil
}

func (r *DictionaryIncrementRequest) getResponse() interface{} {
	return r.response
}
