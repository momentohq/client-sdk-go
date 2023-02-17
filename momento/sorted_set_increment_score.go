package momento

import (
	"context"
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"github.com/momentohq/client-sdk-go/utils"
)

//////// Response

type SortedSetIncrementScoreResponse interface {
	isSortedSetIncrementResponse()
}
type SortedSetIncrementScoreSuccess struct {
	Value float64
}

func (SortedSetIncrementScoreSuccess) isSortedSetIncrementResponse() {}

////// Request

type SortedSetIncrementScoreRequest struct {
	CacheName     string
	SetName       string
	ElementName   Value
	Amount        float64
	CollectionTTL utils.CollectionTTL

	grpcRequest  *pb.XSortedSetIncrementRequest
	grpcResponse *pb.XSortedSetIncrementResponse
	response     SortedSetIncrementScoreResponse
}

func (r *SortedSetIncrementScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetIncrementScoreRequest) requestName() string { return "Sorted set increment" }

func (r *SortedSetIncrementScoreRequest) ttl() time.Duration { return r.CollectionTTL.Ttl }

func (r *SortedSetIncrementScoreRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var ttlMillis uint64
	if ttlMillis, err = prepareTTL(r, client.defaultTtl); err != nil {
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
		ElementName:     r.ElementName.asBytes(),
		Amount:          r.Amount,
		TtlMilliseconds: ttlMillis,
		RefreshTtl:      r.CollectionTTL.RefreshTtl,
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
	r.response = &SortedSetIncrementScoreSuccess{
		Value: r.grpcResponse.Value,
	}
	return nil
}
