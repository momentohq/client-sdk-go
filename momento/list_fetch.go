package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListFetchRequest struct {
	CacheName  string
	ListName   string
	StartIndex *int32
	EndIndex   *int32

	grpcRequest  *pb.XListFetchRequest
	grpcResponse *pb.XListFetchResponse
	response     responses.ListFetchResponse
}

func (r *ListFetchRequest) cacheName() string { return r.CacheName }

func (r *ListFetchRequest) requestName() string { return "ListFetch" }

func (r *ListFetchRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	grpcRequest := &pb.XListFetchRequest{
		ListName:   []byte(r.ListName),
		StartIndex: &pb.XListFetchRequest_UnboundedStart{},
		EndIndex:   &pb.XListFetchRequest_UnboundedEnd{},
	}

	if r.StartIndex != nil {
		grpcRequest.StartIndex = &pb.XListFetchRequest_InclusiveStart{
			InclusiveStart: *r.StartIndex,
		}
	}

	if r.EndIndex != nil {
		grpcRequest.EndIndex = &pb.XListFetchRequest_ExclusiveEnd{
			ExclusiveEnd: *r.EndIndex,
		}
	}

	r.grpcRequest = grpcRequest
	return nil
}

func (r *ListFetchRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListFetch(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	r.grpcResponse = resp

	return resp, nil, nil
}

func (r *ListFetchRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.List.(type) {
	case *pb.XListFetchResponse_Found:
		r.response = responses.NewListFetchHit(rtype.Found.Values)
	case *pb.XListFetchResponse_Missing:
		r.response = &responses.ListFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func (r *ListFetchRequest) getResponse() interface{} { return r.response }
