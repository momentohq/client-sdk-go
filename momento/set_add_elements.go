package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

type SetAddElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []Value
	Ttl       *utils.CollectionTtl

}

func (r *SetAddElementsRequest) cacheName() string { return r.CacheName }

func (r *SetAddElementsRequest) ttl() time.Duration {
	return r.Ttl.Ttl
}

func (r *SetAddElementsRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *SetAddElementsRequest) requestName() string { return "SetAddElements" }

func (r *SetAddElementsRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	elements, err := momentoValuesToPrimitiveByteList(r.Elements)
	if err != nil {
		return nil, err
	}

	grpcRequest := &pb.XSetUnionRequest{
		SetName:         []byte(r.SetName),
		Elements:        elements,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}

	return grpcRequest, nil
}

func (r *SetAddElementsRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetUnion(requestMetadata, grpcRequest.(*pb.XSetUnionRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetAddElementsRequest) interpretGrpcResponse(_ interface{}) (interface{}, error) {
	return &responses.SetAddElementsSuccess{}, nil
}
