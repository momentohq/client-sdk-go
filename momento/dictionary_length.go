package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionaryLengthRequest struct {
	CacheName      string
	DictionaryName string
}

func (r *DictionaryLengthRequest) cacheName() string { return r.CacheName }

func (r *DictionaryLengthRequest) requestName() string { return "DictionaryLength" }

func (r *DictionaryLengthRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	if _, err := prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XDictionaryLengthRequest{
		DictionaryName: []byte(r.DictionaryName),
	}

	return grpcRequest, nil
}

func (r *DictionaryLengthRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionaryLength(requestMetadata, grpcRequest.(*pb.XDictionaryLengthRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *DictionaryLengthRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XDictionaryLengthResponse)
	switch rtype := myResp.Dictionary.(type) {
	case *pb.XDictionaryLengthResponse_Found:
		return responses.NewDictionaryLengthHit(rtype.Found.Length), nil
	case *pb.XDictionaryLengthResponse_Missing:
		return &responses.DictionaryLengthMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}

func (c DictionaryLengthRequest) GetRequestName() string {
	return "DictionaryLengthRequest"
}
