package momento

import (
	"context"
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"github.com/momentohq/client-sdk-go/utils"
)

type SortedSetIncrementScoreRequest struct {
	CacheName string
	SetName   string
	Value     Value
	Amount    float64
	Ttl       *utils.CollectionTtl
}

func (r *SortedSetIncrementScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetIncrementScoreRequest) requestName() string { return "SortedSetIncrementScore" }

func (r *SortedSetIncrementScoreRequest) value() Value { return r.Value }

func (r *SortedSetIncrementScoreRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *SortedSetIncrementScoreRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *SortedSetIncrementScoreRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	if r.Amount == 0 {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Amount must be given and cannot be 0",
			errors.New("invalid argument"),
		)
	}

	grpcRequest := &pb.XSortedSetIncrementRequest{
		SetName:         []byte(r.SetName),
		Value:           value,
		Amount:          r.Amount,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}
	return grpcRequest, nil
}

func (r *SortedSetIncrementScoreRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetIncrement(requestMetadata, grpcRequest.(*pb.XSortedSetIncrementRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetIncrementScoreRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSortedSetIncrementResponse)
	return responses.SortedSetIncrementScoreSuccess(myResp.Score), nil
}
