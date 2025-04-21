package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionaryGetFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Fields         []Value

	grpcResponse *pb.XDictionaryGetResponse
}

func (r *DictionaryGetFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionaryGetFieldsRequest) fields() []Value { return r.Fields }

func (r *DictionaryGetFieldsRequest) requestName() string { return "DictionaryGetFields" }

func (r *DictionaryGetFieldsRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return nil, err
	}

	var fields [][]byte
	if fields, err = prepareFields(r); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XDictionaryGetRequest{
		DictionaryName: []byte(r.DictionaryName),
		Fields:         fields,
	}

	return grpcRequest, nil
}

func (r *DictionaryGetFieldsRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionaryGet(requestMetadata, grpcRequest.(*pb.XDictionaryGetRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *DictionaryGetFieldsRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	r.grpcResponse = resp.(*pb.XDictionaryGetResponse)
	switch rtype := r.grpcResponse.Dictionary.(type) {
	case *pb.XDictionaryGetResponse_Missing:
		return &responses.DictionaryGetFieldsMiss{}, nil
	case *pb.XDictionaryGetResponse_Found:
		var responsesToReturn []responses.DictionaryGetFieldResponse
		var fields [][]byte
		for idx, val := range rtype.Found.Items {
			field := r.Fields[idx].asBytes()
			if val.Result == pb.ECacheResult_Hit {
				responsesToReturn = append(responsesToReturn, responses.NewDictionaryGetFieldHit(field, val.CacheBody))
			} else if val.Result == pb.ECacheResult_Miss {
				responsesToReturn = append(responsesToReturn, responses.NewDictionaryGetFieldMiss(field))
			} else {
				responsesToReturn = append(responsesToReturn, nil)
			}
			fields = append(fields, field)
		}
		return responses.NewDictionaryGetFieldsHit(fields, rtype.Found.Items, responsesToReturn), nil
	default:
		return nil, errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
}

func (c DictionaryGetFieldsRequest) GetRequestName() string {
	return "DictionaryGetFieldsRequest"
}
