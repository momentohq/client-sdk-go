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

func (r *SetAddElementsRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetUnion(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *SetAddElementsRequest) interpretGrpcResponse() error {
	r.response = &responses.SetAddElementsSuccess{}
	return nil
}

func (r *SetAddElementsRequest) getResponse() map[string]string {
	return getMomentoResponseData(r.response)
}
