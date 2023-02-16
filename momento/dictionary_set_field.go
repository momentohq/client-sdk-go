package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionarySetFieldResponse

type DictionarySetFieldResponse interface {
	isDictionarySetFieldResponse()
}

type DictionarySetFieldSuccess struct{}

func (DictionarySetFieldSuccess) isDictionarySetFieldResponse() {}

// DictionarySetFieldRequest

type DictionarySetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Bytes
	Value          Bytes
	CollectionTTL  utils.CollectionTTL

	grpcRequest  *pb.XDictionarySetRequest
	grpcResponse *pb.XDictionarySetResponse
	response     DictionarySetFieldResponse
}

func (r *DictionarySetFieldRequest) cacheName() string { return r.CacheName }

func (r *DictionarySetFieldRequest) field() Bytes { return r.Field }

func (r *DictionarySetFieldRequest) value() Bytes { return r.Value }

func (r *DictionarySetFieldRequest) ttl() time.Duration { return r.CollectionTTL.Ttl }

func (r *DictionarySetFieldRequest) requestName() string { return "DictionarySetField" }

func (r *DictionarySetFieldRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var field []byte
	if field, err = prepareField(r); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	var items []*pb.XDictionaryFieldValuePair
	items = append(items, &pb.XDictionaryFieldValuePair{
		Field: field,
		Value: value,
	})

	r.grpcRequest = &pb.XDictionarySetRequest{
		DictionaryName:  []byte(r.DictionaryName),
		Items:           items,
		TtlMilliseconds: ttl,
		RefreshTtl:      r.CollectionTTL.RefreshTtl,
	}

	return nil
}

func (r *DictionarySetFieldRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionarySet(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionarySetFieldRequest) interpretGrpcResponse() error {
	r.response = &DictionarySetFieldSuccess{}
	return nil
}
