package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionaryLengthRequest struct {
	CacheName      string
	DictionaryName string

	grpcRequest  *pb.XDictionaryLengthRequest
	grpcResponse *pb.XDictionaryLengthResponse
	response     responses.DictionaryLengthResponse
}

func (r *DictionaryLengthRequest) cacheName() string { return r.CacheName }

func (r *DictionaryLengthRequest) requestName() string { return "DictionaryLength" }

func (r *DictionaryLengthRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XDictionaryLengthRequest{
		DictionaryName: []byte(r.DictionaryName),
	}

	return nil
}

func (r *DictionaryLengthRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionaryLength(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionaryLengthRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	switch rtype := resp.Dictionary.(type) {
	case *pb.XDictionaryLengthResponse_Found:
		r.response = responses.NewDictionaryLengthHit(rtype.Found.Length)
	case *pb.XDictionaryLengthResponse_Missing:
		r.response = &responses.DictionaryLengthMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
