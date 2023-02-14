package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// ListFetchResponse

type ListFetchResponse interface {
	isListFetchResponse()
}

type ListFetchHit struct {
	value       [][]byte
	stringValue []string
}

func (ListFetchHit) isListFetchResponse() {}

func (resp ListFetchHit) ValueListByte() [][]byte {
	return resp.value
}

func (resp ListFetchHit) ValueListString() []string {
	if resp.stringValue == nil {
		for _, element := range resp.value {
			resp.stringValue = append(resp.stringValue, string(element))
		}
	}
	return resp.stringValue
}

func (resp ListFetchHit) ValueList() []string {
	return resp.ValueListString()
}

type ListFetchMiss struct{}

func (ListFetchMiss) isListFetchResponse() {}

// ListFetchRequest

type ListFetchRequest struct {
	CacheName string
	ListName  string

	grpcRequest  *pb.XListFetchRequest
	grpcResponse *pb.XListFetchResponse
	response     ListFetchResponse
}

func (r ListFetchRequest) cacheName() string { return r.CacheName }

func (r ListFetchRequest) requestName() string { return "ListFetch" }

func (r *ListFetchRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListFetchRequest{
		ListName: []byte(r.ListName),
	}

	return nil
}

func (r *ListFetchRequest) makeGrpcRequest(client scsDataClient, ctx context.Context) (grpcResponse, error) {
	resp, err := client.grpcClient.ListFetch(ctx, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *ListFetchRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.List.(type) {
	case *pb.XListFetchResponse_Found:
		r.response = ListFetchHit{value: rtype.Found.Values}
	case *pb.XListFetchResponse_Missing:
		r.response = ListFetchMiss{}
	default:
		return errUnexpectedGrpcResponse
	}
	return nil
}
