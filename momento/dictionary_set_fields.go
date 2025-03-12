package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionarySetFieldsRequest represents a request to store multiple elements in a Dictionary.
//
//	Use momento.DictionaryElementsFromMap to help construct the Request from a map object.
type DictionarySetFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Elements       []DictionaryElement
	Ttl            *utils.CollectionTtl

	grpcRequest  *pb.XDictionarySetRequest

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

func (r *DictionarySetFieldsRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionarySet(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *DictionarySetFieldsRequest) interpretGrpcResponse(_ interface{}) error {
	r.response = &responses.DictionarySetFieldsSuccess{}
	return nil
}

func (r *DictionarySetFieldsRequest) getResponse() interface{} {
	return r.response
}
