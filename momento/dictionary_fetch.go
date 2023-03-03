package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionaryFetchRequest struct {
	CacheName      string
	DictionaryName string

	grpcRequest  *pb.XDictionaryFetchRequest
	grpcResponse *pb.XDictionaryFetchResponse
	response     responses.DictionaryFetchResponse
}

func (r *DictionaryFetchRequest) cacheName() string { return r.CacheName }

func (r *DictionaryFetchRequest) requestName() string { return "DictionaryFetch" }

func (r *DictionaryFetchRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XDictionaryFetchRequest{DictionaryName: []byte(r.DictionaryName)}

	return nil
}

func (r *DictionaryFetchRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionaryFetch(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionaryFetchRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.Dictionary.(type) {
	case *pb.XDictionaryFetchResponse_Found:
		elements := make(map[string][]byte)
		for _, i := range rtype.Found.Items {
			elements[(string(i.Field))] = i.Value
		}
		r.response = responses.NewDictionaryFetchHit(elements)
	case *pb.XDictionaryFetchResponse_Missing:
		r.response = &responses.DictionaryFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
