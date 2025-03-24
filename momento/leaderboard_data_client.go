package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type leaderboardDataClient struct {
	requestTimeout         time.Duration
	leaderboardGrpcManager *grpcmanagers.LeaderboardGrpcManager
	leaderboardClient      pb.LeaderboardClient
}

func newLeaderboardDataClient(request *models.LeaderboardClientRequest) (*leaderboardDataClient, momentoerrors.MomentoSvcErr) {
	grpcManager, err := grpcmanagers.NewLeaderboardGrpcManager(&models.LeaderboardGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  request.Configuration.GetTransportStrategy().GetGrpcConfig(),
	})
	if err != nil {
		return nil, err
	}
	return &leaderboardDataClient{
		requestTimeout:         request.Configuration.GetClientSideTimeout(),
		leaderboardGrpcManager: grpcManager,
		leaderboardClient:      pb.NewLeaderboardClient(grpcManager.Conn),
	}, nil
}

func (client *leaderboardDataClient) close() momentoerrors.MomentoSvcErr {
	return client.leaderboardGrpcManager.Close()
}

func (client *leaderboardDataClient) delete(ctx context.Context, request *LeaderboardInternalDeleteRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateLeaderboardMetadata(ctx, request.CacheName)

	var header, trailer metadata.MD
	_, err := client.leaderboardClient.DeleteLeaderboard(requestMetadata, &pb.XDeleteLeaderboardRequest{
		Leaderboard: request.LeaderboardName,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return nil
}

func (client *leaderboardDataClient) fetchByRank(ctx context.Context, request *LeaderboardInternalFetchByRankRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateLeaderboardMetadata(ctx, request.CacheName)

	rankRange := &pb.XRankRange{
		StartInclusive: request.StartRank,
		EndExclusive:   request.EndRank,
	}

	leaderboardOrder := pb.XOrder_ASCENDING
	if request.Order != nil && *request.Order == DESCENDING {
		leaderboardOrder = pb.XOrder_DESCENDING
	}

	var header, trailer metadata.MD
	result, err := client.leaderboardClient.GetByRank(requestMetadata, &pb.XGetByRankRequest{
		Leaderboard: request.LeaderboardName,
		RankRange:   rankRange,
		Order:       leaderboardOrder,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return result.Elements, nil
}

func (client *leaderboardDataClient) fetchByScore(ctx context.Context, request *LeaderboardInternalFetchByScoreRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateLeaderboardMetadata(ctx, request.CacheName)

	scoreRange := &pb.XScoreRange{}

	if request.MinScore == nil {
		scoreRange.Min = &pb.XScoreRange_UnboundedMin{UnboundedMin: &pb.XUnbounded{}}
	} else {
		scoreRange.Min = &pb.XScoreRange_MinInclusive{MinInclusive: *request.MinScore}
	}

	if request.MaxScore == nil {
		scoreRange.Max = &pb.XScoreRange_UnboundedMax{UnboundedMax: &pb.XUnbounded{}}
	} else {
		scoreRange.Max = &pb.XScoreRange_MaxExclusive{MaxExclusive: *request.MaxScore}
	}

	leaderboardOrder := pb.XOrder_ASCENDING
	if request.Order != nil && *request.Order == DESCENDING {
		leaderboardOrder = pb.XOrder_DESCENDING
	}

	offset := uint32(0)
	if request.Offset != nil {
		offset = *request.Offset
	}

	count := uint32(8192)
	if request.Count != nil {
		count = *request.Count
	}

	var header, trailer metadata.MD
	result, err := client.leaderboardClient.GetByScore(requestMetadata, &pb.XGetByScoreRequest{
		Leaderboard:   request.LeaderboardName,
		ScoreRange:    scoreRange,
		Offset:        offset,
		LimitElements: count,
		Order:         leaderboardOrder,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return result.Elements, nil
}

func (client *leaderboardDataClient) getRank(ctx context.Context, request *LeaderboardInternalGetRankRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	leaderboardOrder := pb.XOrder_ASCENDING
	if request.Order != nil && *request.Order == DESCENDING {
		leaderboardOrder = pb.XOrder_DESCENDING
	}

	requestMetadata := internal.CreateLeaderboardMetadata(ctx, request.CacheName)

	var header, trailer metadata.MD
	result, err := client.leaderboardClient.GetRank(requestMetadata, &pb.XGetRankRequest{
		Leaderboard: request.LeaderboardName,
		Ids:         request.Ids,
		Order:       leaderboardOrder,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return result.Elements, nil
}

func (client *leaderboardDataClient) length(ctx context.Context, request *LeaderboardInternalLengthRequest) (uint32, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateLeaderboardMetadata(ctx, request.CacheName)

	var header, trailer metadata.MD
	result, err := client.leaderboardClient.GetLeaderboardLength(requestMetadata, &pb.XGetLeaderboardLengthRequest{
		Leaderboard: request.LeaderboardName,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return 0, momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return result.Count, nil
}

func (client *leaderboardDataClient) removeElements(ctx context.Context, request *LeaderboardInternalRemoveElementsRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateLeaderboardMetadata(ctx, request.CacheName)

	var header, trailer metadata.MD
	_, err := client.leaderboardClient.RemoveElements(requestMetadata, &pb.XRemoveElementsRequest{
		Leaderboard: request.LeaderboardName,
		Ids:         request.Ids,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return nil
}

func (client *leaderboardDataClient) upsert(ctx context.Context, request *LeaderboardInternalUpsertRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateLeaderboardMetadata(ctx, request.CacheName)

	var header, trailer metadata.MD
	_, err := client.leaderboardClient.UpsertElements(requestMetadata, &pb.XUpsertElementsRequest{
		Leaderboard: request.LeaderboardName,
		Elements:    leaderboardUpsertElementToGrpc(request.Elements),
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return nil
}

func leaderboardUpsertElementToGrpc(elements []LeaderboardUpsertElement) []*pb.XElement {
	var grpcElements []*pb.XElement
	for _, element := range elements {
		grpcElements = append(grpcElements, &pb.XElement{
			Id:    element.Id,
			Score: element.Score,
		})
	}
	return grpcElements
}
