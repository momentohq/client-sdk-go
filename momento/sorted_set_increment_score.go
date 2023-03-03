package momento

import (
	"context"
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"github.com/momentohq/client-sdk-go/utils"
)

type SortedSetIncrementScoreRequest struct {
	CacheName    string
	SetName      string
	ElementValue Value
	Amount       float64
	Ttl          *utils.CollectionTtl

	grpcRequest  *pb.XSortedSetIncrementRequest
	grpcResponse *pb.XSortedSetIncrementResponse
	response     responses.SortedSetIncrementScoreResponse
}

func (r *SortedSetIncrementScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetIncrementScoreRequest) requestName() string { return "Sorted set increment" }

func (r *SortedSetIncrementScoreRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *SortedSetIncrementScoreRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *SortedSetIncrementScoreRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareElementValue(r.ElementValue); err != nil {
		return err
	}

	if r.Amount == 0 {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Amount must be given and cannot be 0",
			errors.New("invalid argument"),
		)
	}

	r.grpcRequest = &pb.XSortedSetIncrementRequest{
		SetName:         []byte(r.SetName),
		Value:           value,
		Amount:          r.Amount,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}
	return nil
}

func (r *SortedSetIncrementScoreRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetIncrement(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SortedSetIncrementScoreRequest) interpretGrpcResponse() error {
	r.response = &responses.SortedSetIncrementScoreSuccess{
		Value: r.grpcResponse.Score,
	}
	return nil
}
