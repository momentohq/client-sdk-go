package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionarySetFieldsResponse

type DictionarySetFieldsResponse interface {
	isDictionarySetFieldsResponse()
}

type DictionarySetFieldsSuccess struct{}

func (DictionarySetFieldsSuccess) isDictionarySetFieldsResponse() {}

// DictionarySetFieldsRequest

type DictionarySetFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Items          map[string]Value
	CollectionTTL  utils.CollectionTTL

	grpcRequest  *pb.XDictionarySetRequest
	grpcResponse *pb.XDictionarySetResponse
	response     DictionarySetFieldsResponse
}

func (r *DictionarySetFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionarySetFieldsRequest) items() map[string]Value { return r.Items }

func (r *DictionarySetFieldsRequest) ttl() time.Duration { return r.CollectionTTL.Ttl }

func (r *DictionarySetFieldsRequest) requestName() string { return "DictionarySetFields" }

func (r *DictionarySetFieldsRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var items map[string][]byte
	if items, err = prepareItems(r); err != nil {
		return err
	}

	var pbItems []*pb.XDictionaryFieldValuePair
	for k, v := range items {
		pbItems = append(pbItems, &pb.XDictionaryFieldValuePair{
			Field: []byte(k),
			Value: v,
		})
	}

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XDictionarySetRequest{
		DictionaryName:  []byte(r.DictionaryName),
		Items:           pbItems,
		TtlMilliseconds: ttl,
		RefreshTtl:      r.CollectionTTL.RefreshTtl,
	}

	return nil
}

func (r *DictionarySetFieldsRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionarySet(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionarySetFieldsRequest) interpretGrpcResponse() error {
	r.response = &DictionarySetFieldsSuccess{}
	return nil
}
