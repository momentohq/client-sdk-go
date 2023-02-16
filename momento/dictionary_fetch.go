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
	items             map[Bytes]Bytes
	itemsStringString map[string]string
	itemsStringByte   map[string]Bytes
}

func (DictionaryFetchHit) isDictionaryFetchResponse() {}

func (resp DictionaryFetchHit) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

func (resp DictionaryFetchHit) ValueMapStringString() map[string]string {
	if resp.itemsStringString == nil {
		resp.itemsStringString = make(map[string]string)
		for k, v := range resp.items {
			resp.itemsStringString[string(k.AsBytes())] = string(v.AsBytes())
		}
	}
	return resp.itemsStringString
}

func (resp DictionaryFetchHit) ValueMapStringByte() map[string]Bytes {
	if resp.itemsStringByte == nil {
		resp.itemsStringByte = make(map[string]Bytes)
		for k, v := range resp.items {
			resp.itemsStringByte[string(k.AsBytes())] = v
		}
	}
	return resp.itemsStringByte
}

func (resp DictionaryFetchHit) ValueMapByteByte() map[Bytes]Bytes {
	return resp.items
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

func (r *DictionaryFetchRequest) initGrpcRequest(client scsDataClient) error {
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
		// TODO: refactor to utility func
		itemsAsBytes := make(map[Bytes]Bytes)
		for _, i := range rtype.Found.Items {
			itemsAsBytes[StringBytes{Text: string(i.Field)}] = RawBytes{i.Value}
		}
		r.response = &DictionaryFetchHit{items: itemsAsBytes}
	case *pb.XDictionaryFetchResponse_Missing:
		r.response = &DictionaryFetchMiss{}
	default:
		return errUnexpectedGrpcResponse
	}
	return nil
}
