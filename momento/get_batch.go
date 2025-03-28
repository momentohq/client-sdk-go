package momento

import (
	"context"
	"io"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc/metadata"
)

type GetBatchRequest struct {
	CacheName string
	Keys      []Value

	grpcStream pb.Scs_GetBatchClient
	byteKeys   [][]byte
}

func (r *GetBatchRequest) cacheName() string { return r.CacheName }

func (r *GetBatchRequest) keys() []Value { return r.Keys }

func (r *GetBatchRequest) requestName() string { return "GetBatch" }

func (r *GetBatchRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.CacheName, "Cache name"); err != nil {
		return nil, err
	}

	// For each key, prepare a GetRequest
	var getRequests []*pb.XGetRequest
	for _, key := range r.Keys {
		var byteKey = key.asBytes()
		r.byteKeys = append(r.byteKeys, byteKey)
		getRequests = append(getRequests, &pb.XGetRequest{
			CacheKey: byteKey,
		})
	}

	grpcRequest := &pb.XGetBatchRequest{
		Items: getRequests,
	}

	return grpcRequest, nil
}

func (r *GetBatchRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	var responseMetadata []metadata.MD
	resp, err := client.grpcClient.GetBatch(requestMetadata, grpcRequest.(*pb.XGetBatchRequest))
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

func (r *GetBatchRequest) interpretGrpcResponse(_ interface{}) (interface{}, error) {
	var getResponses []responses.GetResponse
	for {
		resp, err := r.grpcStream.Recv()
		if err == io.EOF {
			break
		} else if err == nil {
			switch resp.Result {
			case pb.ECacheResult_Hit:
				var getHit = responses.NewGetHit(resp.CacheBody)
				getResponses = append(getResponses, getHit)
			case pb.ECacheResult_Miss:
				getResponses = append(getResponses, &responses.GetMiss{})
			default:
				return nil, momentoerrors.ConvertSvcErr(err)
			}
		} else {
			return nil, momentoerrors.ConvertSvcErr(err)
		}
	}

	return *responses.NewGetBatchSuccess(getResponses, r.byteKeys), nil
}
