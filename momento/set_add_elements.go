package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

// SetAddElementsResponse

type SetAddElementsResponse interface {
	isSetAddElementResponse()
}

type SetAddElementsSuccess struct{}

func (SetAddElementsSuccess) isSetAddElementResponse() {}

// SetAddElementRequest

type SetAddElementsRequest struct {
	CacheName     string
	SetName       string
	Elements      []Value
	CollectionTTL utils.CollectionTTL

	grpcRequest  *pb.XSetUnionRequest
	grpcResponse *pb.XSetUnionResponse
	response     SetAddElementsResponse
}

func (r *SetAddElementsRequest) cacheName() string { return r.CacheName }

func (r *SetAddElementsRequest) ttl() time.Duration {
	return processCollectionTtl(r.CollectionTTL)
}

func (r *SetAddElementsRequest) requestName() string { return "SetAddElements" }

func (r *SetAddElementsRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XSetUnionRequest{
		SetName:         []byte(r.SetName),
		Elements:        momentoValuesToPrimitiveByteList(r.Elements),
		TtlMilliseconds: ttl,
		RefreshTtl:      r.CollectionTTL.RefreshTtl,
	}

	return nil
}

func (r *SetAddElementsRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetUnion(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetAddElementsRequest) interpretGrpcResponse() error {
	r.response = &SetAddElementsSuccess{}
	return nil
}
