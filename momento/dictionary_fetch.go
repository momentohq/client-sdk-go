package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionaryFetchRequest struct {
	CacheName      string
	DictionaryName string
}

func (r *DictionaryFetchRequest) cacheName() string { return r.CacheName }

func (r *DictionaryFetchRequest) requestName() string { return "DictionaryFetch" }

func (r *DictionaryFetchRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XDictionaryFetchRequest{DictionaryName: []byte(r.DictionaryName)}

	return grpcRequest, nil
}

func (r *DictionaryFetchRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionaryFetch(requestMetadata, grpcRequest.(*pb.XDictionaryFetchRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *DictionaryFetchRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XDictionaryFetchResponse)
	switch rtype := myResp.Dictionary.(type) {
	case *pb.XDictionaryFetchResponse_Found:
		elements := make(map[string][]byte)
		for _, i := range rtype.Found.Items {
			elements[(string(i.Field))] = i.Value
		}
		return responses.NewDictionaryFetchHit(elements), nil
	case *pb.XDictionaryFetchResponse_Missing:
		return &responses.DictionaryFetchMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}

func (c DictionaryFetchRequest) GetRequestName() string {
	return "DictionaryFetchRequest"
}
