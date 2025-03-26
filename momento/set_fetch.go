package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetFetchRequest struct {
	CacheName string
	SetName   string
}

func (r *SetFetchRequest) cacheName() string { return r.CacheName }

func (r *SetFetchRequest) requestName() string { return "SetFetch" }

func (r *SetFetchRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XSetFetchRequest{SetName: []byte(r.SetName)}

	return grpcRequest, nil
}

func (r *SetFetchRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetFetch(requestMetadata, grpcRequest.(*pb.XSetFetchRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetFetchRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSetFetchResponse)
	switch rtype := myResp.Set.(type) {
	case *pb.XSetFetchResponse_Found:
		return responses.NewSetFetchHit(rtype.Found.Elements), nil
	case *pb.XSetFetchResponse_Missing:
		return &responses.SetFetchMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
