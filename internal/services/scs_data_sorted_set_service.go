package services

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc/metadata"
)

func (client *ScsDataClient) SortedSetPut(ctx context.Context, request *models.SortedSetPutRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	_, err := client.grpcClient.SortedSetPut(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSortedSetPutRequest{
			SetName:         request.SetName,
			Elements:        sortedSetModelToGrpc(request.Elements),
			TtlMilliseconds: collectionTtlOrDefaultMilliseconds(request.CollectionTTL, client.defaultTtl),
			RefreshTtl:      request.CollectionTTL.RefreshTtl,
		},
	)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsDataClient) SortedSetFetch(ctx context.Context, request *models.SortedSetFetchRequest) (models.SortedSetFetchResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestToMake := &pb.XSortedSetFetchRequest{
		SetName: request.SetName,
		Order:   pb.XSortedSetFetchRequest_Order(request.Order),
	}
	switch r := request.NumberOfResults.(type) {
	case *models.FetchAllElements:
		requestToMake.NumResults = &pb.XSortedSetFetchRequest_All{}
	case *models.FetchLimitedElements:
		requestToMake.NumResults = &pb.XSortedSetFetchRequest_Limit{
			Limit: &pb.XSortedSetFetchRequest_XLimit{
				Limit: r.Limit,
			},
		}
	}
	resp, err := client.grpcClient.SortedSetFetch(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		requestToMake,
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	// Convert from grpc struct to internal struct
	switch r := resp.SortedSet.(type) {
	case *pb.XSortedSetFetchResponse_Found:
		return &models.SortedSetFetchFound{
			Elements: sortedSetGrpcElementToModel(r.Found.GetElements()),
		}, nil
	case *pb.XSortedSetFetchResponse_Missing:
		return &models.SortedSetFetchMissing{}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.ClientSdkError,
			fmt.Sprintf("Unknown response type for sortedset fetch response %+v", r),
			nil,
		)
	}
}
func (client *ScsDataClient) SortedSetGetScore(ctx context.Context, request *models.SortedSetGetScoreRequest) (models.SortedSetGetScoreResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	resp, err := client.grpcClient.SortedSetGetScore(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSortedSetGetScoreRequest{
			SetName:     request.SetName,
			ElementName: request.ElementNames,
		},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	// Convert from grpc struct to internal struct
	switch r := resp.SortedSet.(type) {
	case *pb.XSortedSetGetScoreResponse_Found:
		return &models.SortedSetGetScoreHit{
			Elements: sortedSetGrpcScoreToModel(r.Found.GetElements()),
		}, nil
	case *pb.XSortedSetGetScoreResponse_Missing:
		return &models.SortedSetGetScoreMiss{}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.ClientSdkError,
			fmt.Sprintf("Unknown response type for sortedset GetScore response %+v", r),
			nil,
		)
	}
}

func (client *ScsDataClient) SortedSetGetRank(ctx context.Context, request *models.SortedSetGetRankRequest) (models.SortedSetGetRankResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	resp, err := client.grpcClient.SortedSetGetRank(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSortedSetGetRankRequest{
			SetName:     request.SetName,
			ElementName: request.ElementName,
		},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	// Convert from grpc struct to internal struct
	switch r := resp.Rank.(type) {
	case *pb.XSortedSetGetRankResponse_ElementRank:
		return &models.SortedSetGetRankHit{
			Rank:   r.ElementRank.Rank,
			Status: models.CacheResult(r.ElementRank.Result),
		}, nil
	case *pb.XSortedSetGetRankResponse_Missing:
		return &models.SortedSetGetRankMiss{}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.ClientSdkError,
			fmt.Sprintf("Unknown response type for sortedset GetRank response %+v", r),
			nil,
		)
	}
}

func (client *ScsDataClient) SortedSetRemove(ctx context.Context, request *models.SortedSetRemoveRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	requestToMake := &pb.XSortedSetRemoveRequest{
		SetName: request.SetName,
	}
	switch r := request.ElementsToRemove.(type) {
	case *models.RemoveAllElements:
		requestToMake.RemoveElements = &pb.XSortedSetRemoveRequest_All{}
	case *models.RemoveSomeElements:
		requestToMake.RemoveElements = &pb.XSortedSetRemoveRequest_Some{
			Some: &pb.XSortedSetRemoveRequest_XSome{
				ElementName: r.ElementsToRemove,
			},
		}
	}
	_, err := client.grpcClient.SortedSetRemove(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		requestToMake,
	)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func sortedSetGrpcElementToModel(grpcSetElements []*pb.XSortedSetElement) []*models.SortedSetElement {
	var returnList []*models.SortedSetElement
	for i := range grpcSetElements {
		returnList = append(returnList, &models.SortedSetElement{
			Name:  grpcSetElements[i].Name,
			Score: grpcSetElements[i].Score,
		})
	}
	return returnList
}

func sortedSetGrpcScoreToModel(grpcSetElements []*pb.XSortedSetGetScoreResponse_XSortedSetGetScoreResponsePart) []*models.SortedSetScore {
	var returnList []*models.SortedSetScore
	for i := range grpcSetElements {
		returnList = append(returnList, &models.SortedSetScore{
			Result: models.CacheResult(grpcSetElements[i].Result),
			Score:  grpcSetElements[i].Score,
		})
	}
	return returnList
}

func sortedSetModelToGrpc(modelSetElements []*models.SortedSetElement) []*pb.XSortedSetElement {
	var returnList []*pb.XSortedSetElement
	for _, el := range modelSetElements {
		returnList = append(returnList, &pb.XSortedSetElement{
			Name:  el.Name,
			Score: el.Score,
		})
	}
	return returnList
}
