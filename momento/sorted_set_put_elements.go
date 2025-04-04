package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetPutElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []SortedSetElement
	Ttl       *utils.CollectionTtl
}

func (r *SortedSetPutElementsRequest) cacheName() string { return r.CacheName }

func (r *SortedSetPutElementsRequest) requestName() string { return "SortedSetPutElements" }

func (r *SortedSetPutElementsRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *SortedSetPutElementsRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *SortedSetPutElementsRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return nil, err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	elements, err := convertSortedSetElementsToGrpc(r.Elements)
	if err != nil {
		return nil, err
	}

	grpcRequest := &pb.XSortedSetPutRequest{
		SetName:         []byte(r.SetName),
		Elements:        elements,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}
	return grpcRequest, nil
}

func (r *SortedSetPutElementsRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SortedSetPut(requestMetadata, grpcRequest.(*pb.XSortedSetPutRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SortedSetPutElementsRequest) interpretGrpcResponse(_ interface{}) (interface{}, error) {
	return &responses.SortedSetPutElementsSuccess{}, nil
}

func convertSortedSetElementsToGrpc(modelSetElements []SortedSetElement) ([]*pb.XSortedSetElement, error) {
	if modelSetElements == nil {
		return nil, buildError(
			momentoerrors.InvalidArgumentError, "elements cannot be nil", nil,
		)
	}
	var returnList []*pb.XSortedSetElement
	for _, el := range modelSetElements {
		if el.Value == nil {
			return nil, buildError(
				momentoerrors.InvalidArgumentError, "element value cannot be nil", nil,
			)
		}
		returnList = append(returnList, &pb.XSortedSetElement{
			Value: el.Value.asBytes(),
			Score: el.Score,
		})
	}
	return returnList, nil
}
