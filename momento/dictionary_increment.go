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
}

func (r *DictionaryIncrementRequest) cacheName() string { return r.CacheName }

func (r *DictionaryIncrementRequest) field() Value { return r.Field }

func (r *DictionaryIncrementRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *DictionaryIncrementRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *DictionaryIncrementRequest) requestName() string { return "DictionaryIncrement" }

func (r *DictionaryIncrementRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return nil, err
	}

	var field []byte
	if field, err = prepareField(r); err != nil {
		return nil, err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	if r.Amount == 0 {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Amount must be given and cannot be 0",
			errors.New("invalid argument"),
		)
	}

	grpcRequest := &pb.XDictionaryIncrementRequest{
		DictionaryName:  []byte(r.DictionaryName),
		Field:           field,
		Amount:          r.Amount,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}

	return grpcRequest, nil
}

func (r *DictionaryIncrementRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionaryIncrement(requestMetadata, grpcRequest.(*pb.XDictionaryIncrementRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *DictionaryIncrementRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XDictionaryIncrementResponse)
	return responses.NewDictionaryIncrementSuccess(myResp.Value), nil
}

func (c DictionaryIncrementRequest) GetRequestName() string {
	return "DictionaryIncrementRequest"
}
