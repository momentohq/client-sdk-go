package momento

import (
	"context"
	"io"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc/metadata"
)

type SetBatchRequest struct {
	CacheName string
	Items     []BatchSetItem
	Ttl       time.Duration

	grpcStream pb.Scs_SetBatchClient
	response   responses.SetBatchResponse
}

func (r *SetBatchRequest) cacheName() string { return r.CacheName }

func (r *SetBatchRequest) ttl() time.Duration { return r.Ttl }

func (r *SetBatchRequest) requestName() string { return "SetBatch" }

func (r *SetBatchRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error
	if _, err = prepareName(r.CacheName, "Cache name"); err != nil {
		return nil, err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return nil, err
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

	grpcRequest := &pb.XSetBatchRequest{
		Items: setRequests,
	}

	return grpcRequest, nil
}

func (r *SetBatchRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	var responseMetadata []metadata.MD
	resp, err := client.grpcClient.SetBatch(requestMetadata, grpcRequest.(*pb.XSetBatchRequest))
	// If there is an error, it's possible resp is nil and we should avoid
	// calling Header() and Trailer() on it to avoid a panic
	if resp != nil {
		header, _ = resp.Header()
		trailer = resp.Trailer()
		responseMetadata = []metadata.MD{header, trailer}
	}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcStream = resp
	// Not sure what to return here, don't think it's even used
	return nil, nil, nil
}

func (r *SetBatchRequest) interpretGrpcResponse(_ interface{}) error {
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

func (r *SetBatchRequest) validateResponseType(resp grpcResponse) error {
	return nil
}
