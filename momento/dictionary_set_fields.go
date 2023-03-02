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
	Elements       map[string]Value
	Ttl            *utils.CollectionTtl

	grpcRequest  *pb.XDictionarySetRequest
	grpcResponse *pb.XDictionarySetResponse
	response     DictionarySetFieldsResponse
}

func (r *DictionarySetFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionarySetFieldsRequest) elements() map[string]Value { return r.Elements }

func (r *DictionarySetFieldsRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *DictionarySetFieldsRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *DictionarySetFieldsRequest) requestName() string { return "DictionarySetFields" }

func (r *DictionarySetFieldsRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var elements map[string][]byte
	if elements, err = prepareElements(r); err != nil {
		return err
	}

	var pbElements []*pb.XDictionaryFieldValuePair
	for k, v := range elements {
		pbElements = append(pbElements, &pb.XDictionaryFieldValuePair{
			Field: []byte(k),
			Value: v,
		})
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XDictionarySetRequest{
		DictionaryName:  []byte(r.DictionaryName),
		Items:           pbElements,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
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
