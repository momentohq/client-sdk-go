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
}

func (r *ListFetchRequest) cacheName() string { return r.CacheName }

func (r *ListFetchRequest) requestName() string { return "ListFetch" }

func (r *ListFetchRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return nil, err
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

	return grpcRequest, nil
}

func (r *ListFetchRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListFetch(requestMetadata, grpcRequest.(*pb.XListFetchRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListFetchRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XListFetchResponse)
	switch rtype := myResp.List.(type) {
	case *pb.XListFetchResponse_Found:
		return responses.NewListFetchHit(rtype.Found.Values), nil
	case *pb.XListFetchResponse_Missing:
		return &responses.ListFetchMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
