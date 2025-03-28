package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetLengthRequest struct {
	CacheName string
	SetName   string
}

func (r *SortedSetLengthRequest) cacheName() string { return r.CacheName }

func (r *SortedSetLengthRequest) requestName() string { return "SortedSetLength" }

func (r *SortedSetLengthRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	if _, err := prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XSortedSetLengthRequest{
		SetName: []byte(r.SetName),
	}

	return grpcRequest, nil
}

func (r *SortedSetLengthRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetLength(requestMetadata, grpcRequest.(*pb.XSortedSetLengthRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetLengthRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSortedSetLengthResponse)
	switch rtype := myResp.SortedSet.(type) {
	case *pb.XSortedSetLengthResponse_Found:
		return responses.NewSortedSetLengthHit(rtype.Found.Length), nil
	case *pb.XSortedSetLengthResponse_Missing:
		return &responses.SortedSetLengthMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
