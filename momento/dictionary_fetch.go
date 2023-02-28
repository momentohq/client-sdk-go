package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionaryFetchResponse

type DictionaryFetchResponse interface {
	isDictionaryFetchResponse()
}

type DictionaryFetchHit struct {
	elementsStringByte   map[string][]byte
	elementsStringString map[string]string
}

func (DictionaryFetchHit) isDictionaryFetchResponse() {}

func (resp DictionaryFetchHit) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

func (resp DictionaryFetchHit) ValueMapStringString() map[string]string {
	if resp.elementsStringString == nil {
		resp.elementsStringString = make(map[string]string)
		for k, v := range resp.elementsStringByte {
			resp.elementsStringString[k] = string(v)
		}
	}
	return resp.elementsStringString
}

func (resp DictionaryFetchHit) ValueMapStringByte() map[string][]byte {
	return resp.elementsStringByte
}

type DictionaryFetchMiss struct{}

func (DictionaryFetchMiss) isDictionaryFetchResponse() {}

// DictionaryFetchRequest

type DictionaryFetchRequest struct {
	CacheName      string
	DictionaryName string

	grpcRequest  *pb.XDictionaryFetchRequest
	grpcResponse *pb.XDictionaryFetchResponse
	response     DictionaryFetchResponse
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
		r.response = &DictionaryFetchHit{elementsStringByte: elements}
	case *pb.XDictionaryFetchResponse_Missing:
		r.response = &DictionaryFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
