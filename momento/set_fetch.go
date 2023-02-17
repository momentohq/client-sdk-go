package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// SetFetchResponse

type SetFetchResponse interface {
	isSetFetchResponse()
}

type SetFetchHit struct {
	elements       [][]byte
	elementsString []string
}

func (SetFetchHit) isSetFetchResponse() {}

func (resp SetFetchHit) ValueSetString() []string {
	if resp.elementsString == nil {
		for _, value := range resp.elements {
			resp.elementsString = append(resp.elementsString, string(value))
		}
	}
	return resp.elementsString
}

func (resp SetFetchHit) ValueSetByte() [][]byte {
	return resp.elements
}

type SetFetchMiss struct{}

func (SetFetchMiss) isSetFetchResponse() {}

// SetFetchRequest

type SetFetchRequest struct {
	CacheName string
	SetName   string

	grpcRequest  *pb.XSetFetchRequest
	grpcResponse *pb.XSetFetchResponse
	response     SetFetchResponse
}

func (r *SetFetchRequest) cacheName() string { return r.CacheName }

func (r *SetFetchRequest) requestName() string { return "SetFetch" }

func (r *SetFetchRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XSetFetchRequest{SetName: []byte(r.SetName)}

	return nil
}

func (r *SetFetchRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetFetch(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetFetchRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.Set.(type) {
	case *pb.XSetFetchResponse_Found:
		r.response = SetFetchHit{
			elements: rtype.Found.Elements,
		}
	case *pb.XSetFetchResponse_Missing:
		r.response = SetFetchMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
