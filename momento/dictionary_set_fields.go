package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionarySetFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Elements       []DictionaryElement
	Ttl            *utils.CollectionTtl

	grpcRequest  *pb.XDictionarySetRequest
	grpcResponse *pb.XDictionarySetResponse
	response     responses.DictionarySetFieldsResponse
}

func (r *DictionarySetFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionarySetFieldsRequest) dictionaryElements() []DictionaryElement { return r.Elements }

func (r *DictionarySetFieldsRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *DictionarySetFieldsRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *DictionarySetFieldsRequest) requestName() string { return "DictionarySetFields" }

func (r *DictionarySetFieldsRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var elements []DictionaryElement
	if elements, err = prepareDictionaryElements(r); err != nil {
		return err
	}

	var pbElements []*pb.XDictionaryFieldValuePair
	for _, v := range elements {
		pbElements = append(pbElements, &pb.XDictionaryFieldValuePair{
			Field: v.Field.asBytes(),
			Value: v.Value.asBytes(),
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
	r.response = &responses.DictionarySetFieldsSuccess{}
	return nil
}
