package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetLengthByScoreRequest struct {
	CacheName string
	SetName   string
	MinScore  *float64
	MaxScore  *float64

	grpcRequest  *pb.XSortedSetLengthByScoreRequest
	grpcResponse *pb.XSortedSetLengthByScoreResponse
	response     responses.SortedSetLengthByScoreResponse
}

func (r *SortedSetLengthByScoreRequest) cacheName() string { return r.CacheName }

func (r *SortedSetLengthByScoreRequest) requestName() string { return "SortedSetLengthByScore" }

func (r *SortedSetLengthByScoreRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	grpc_request := &pb.XSortedSetLengthByScoreRequest{
		SetName: []byte(r.SetName),
	}

	switch r.MaxScore {
	case nil:
		// if no score is provided, we take unbounded or inf by default
		grpc_request.Max = &pb.XSortedSetLengthByScoreRequest_UnboundedMax{}
	default:
		// if a score is provided, it's inclusive by default
		grpc_request.Max = &pb.XSortedSetLengthByScoreRequest_InclusiveMax{
			InclusiveMax: *r.MaxScore,
		}
	}

	switch r.MinScore {
	case nil:
		// if no score is provided, we take unbounded or -inf by default
		grpc_request.Min = &pb.XSortedSetLengthByScoreRequest_UnboundedMin{}
	default:
		// if a score is provided, it's inclusive by default
		grpc_request.Min = &pb.XSortedSetLengthByScoreRequest_InclusiveMin{
			InclusiveMin: *r.MinScore,
		}
	}

	r.grpcRequest = grpc_request

	return nil
}

func (r *SortedSetLengthByScoreRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetLengthByScore(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *SortedSetLengthByScoreRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	switch rtype := resp.SortedSet.(type) {
	case *pb.XSortedSetLengthByScoreResponse_Found:
		r.response = responses.NewSortedSetLengthByScoreHit(rtype.Found.Length)
	case *pb.XSortedSetLengthByScoreResponse_Missing:
		r.response = &responses.SortedSetLengthByScoreMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
