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

	grpcRequest  *pb.XDictionaryGetRequest
	grpcResponse *pb.XDictionaryGetResponse
	response     responses.DictionaryGetFieldsResponse
}

func (r *DictionaryGetFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionaryGetFieldsRequest) fields() []Value { return r.Fields }

func (r *DictionaryGetFieldsRequest) requestName() string { return "DictionaryGetFields" }

func (r *DictionaryGetFieldsRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var fields [][]byte
	if fields, err = prepareFields(r); err != nil {
		return err
	}

	r.grpcRequest = &pb.XDictionaryGetRequest{
		DictionaryName: []byte(r.DictionaryName),
		Fields:         fields,
	}

	return nil
}

func (r *DictionaryGetFieldsRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionaryGet(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *DictionaryGetFieldsRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XDictionaryGetResponse)
	switch rtype := myResp.Dictionary.(type) {
	case *pb.XDictionaryGetResponse_Missing:
		r.response = &responses.DictionaryGetFieldsMiss{}
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
		r.response = responses.NewDictionaryGetFieldsHit(fields, rtype.Found.Items, responsesToReturn)
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}
