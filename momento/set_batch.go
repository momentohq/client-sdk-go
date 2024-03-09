package momento

import (
	"context"
	"io"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type SetBatchRequest struct {
	CacheName string
	Items     []BatchSetItem
	Ttl       time.Duration

	grpcRequest *pb.XSetBatchRequest
	grpcStream  pb.Scs_SetBatchClient
	response    responses.SetBatchResponse
}

func (r *SetBatchRequest) cacheName() string { return r.CacheName }

func (r *SetBatchRequest) ttl() time.Duration { return r.Ttl }

func (r *SetBatchRequest) requestName() string { return "SetBatch" }

func (r *SetBatchRequest) initGrpcRequest(client scsDataClient) error {
	var err error
	if _, err = prepareName(r.CacheName, "Cache name"); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return err
	}

	// For each item, prepare a SetRequest
	var setRequests []*pb.XSetRequest
	for _, item := range r.Items {
		setRequests = append(setRequests, &pb.XSetRequest{
			CacheKey:        item.Key.asBytes(),
			CacheBody:       item.Value.asBytes(),
			TtlMilliseconds: ttl,
		})
	}

	r.grpcRequest = &pb.XSetBatchRequest{
		Items: setRequests,
	}

	return nil
}

func (r *SetBatchRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetBatch(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcStream = resp
	// Not sure what to return here, don't think it's even used
	return nil, nil
}

func (r *SetBatchRequest) interpretGrpcResponse() error {
	var setResponses []responses.SetResponse
	for {
		resp, err := r.grpcStream.Recv()
		if err == io.EOF {
			break
		} else if err == nil {
			switch resp.Result {
			case pb.ECacheResult_Ok:
				setResponses = append(setResponses, &responses.SetSuccess{})
			default:
				return momentoerrors.ConvertSvcErr(err)
			}
		} else {
			return momentoerrors.ConvertSvcErr(err)
		}
	}

	r.response = *responses.NewSetBatchSuccess(setResponses)
	return nil
}
