package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetPopRequest struct {
	CacheName string
	SetName   string
	Count     *uint32

}

func (r *SetPopRequest) cacheName() string { return r.CacheName }

func (r *SetPopRequest) requestName() string { return "SetPop" }

func (r *SetPopRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	var count uint32 = 1
	if r.Count != nil {
		count = uint32(*r.Count)
	}

	grpcRequest := &pb.XSetPopRequest{
		SetName: []byte(r.SetName),
		Count:   count,
	}

	return grpcRequest, nil
}

func (r *SetPopRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetPop(requestMetadata, grpcRequest.(*pb.XSetPopRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetPopRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSetPopResponse)
	switch rtype := myResp.Set.(type) {
	case *pb.XSetPopResponse_Found:
		return responses.NewSetPopHit(rtype.Found.Elements), nil
	case *pb.XSetPopResponse_Missing:
		return &responses.SetPopMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
