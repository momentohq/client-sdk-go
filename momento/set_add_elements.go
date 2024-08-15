package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

type SetAddElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []Value
	Ttl       *utils.CollectionTtl

	grpcRequest  *pb.XSetUnionRequest
	grpcResponse *pb.XSetUnionResponse
	response     responses.SetAddElementsResponse
}

func (r *SetAddElementsRequest) cacheName() string { return r.CacheName }

func (r *SetAddElementsRequest) ttl() time.Duration {
	return r.Ttl.Ttl
}

func (r *SetAddElementsRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *SetAddElementsRequest) requestName() string { return "SetAddElements" }

func (r *SetAddElementsRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	elements, err := momentoValuesToPrimitiveByteList(r.Elements)
	if err != nil {
		return err
	}

	r.grpcRequest = &pb.XSetUnionRequest{
		SetName:         []byte(r.SetName),
		Elements:        elements,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
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
	r.response = &responses.SetAddElementsSuccess{}
	return nil
}

func (r *SetAddElementsRequest) getResponse() interface{} {
	return r.response
}
