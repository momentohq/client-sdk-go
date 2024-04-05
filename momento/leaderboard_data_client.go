package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type LeaderboardDataClient interface {
	Close() momentoerrors.MomentoSvcErr
	Delete(ctx context.Context, request *LeaderboardInternalDeleteRequest) momentoerrors.MomentoSvcErr
	FetchByRank(ctx context.Context, request *LeaderboardInternalFetchByRankRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr)
	FetchByScore(ctx context.Context, request *LeaderboardInternalFetchByScoreRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr)
	GetRank(ctx context.Context, request *LeaderboardInternalGetRankRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr)
	Length(ctx context.Context, request *LeaderboardInternalLengthRequest) (uint32, momentoerrors.MomentoSvcErr)
	RemoveElements(ctx context.Context, request *LeaderboardInternalRemoveElementsRequest) momentoerrors.MomentoSvcErr
	Upsert(ctx context.Context, request *LeaderboardInternalUpsertRequest) momentoerrors.MomentoSvcErr
}

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

func (client *leaderboardDataClient) Close() momentoerrors.MomentoSvcErr {
	return client.leaderboardGrpcManager.Close()
}

func (client *leaderboardDataClient) Delete(ctx context.Context, request *LeaderboardInternalDeleteRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	_, err := client.leaderboardClient.DeleteLeaderboard(ctx, &pb.XDeleteLeaderboardRequest{
		CacheName:   request.CacheName,
		Leaderboard: request.LeaderboardName,
	})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *leaderboardDataClient) FetchByRank(ctx context.Context, request *LeaderboardInternalFetchByRankRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	rankRange := &pb.XRankRange{
		StartInclusive: request.StartRank,
		EndExclusive:   request.EndRank,
	}

	leaderboardOrder := pb.XOrder_ASCENDING
	if request.Order != nil && *request.Order == DESCENDING {
		leaderboardOrder = pb.XOrder_DESCENDING
	}

	result, err := client.leaderboardClient.GetByRank(ctx, &pb.XGetByRankRequest{
		CacheName:   request.CacheName,
		Leaderboard: request.LeaderboardName,
		RankRange:   rankRange,
		Order:       leaderboardOrder,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return result.Elements, nil
}

func (client *leaderboardDataClient) FetchByScore(ctx context.Context, request *LeaderboardInternalFetchByScoreRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

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

	result, err := client.leaderboardClient.GetByScore(ctx, &pb.XGetByScoreRequest{
		CacheName:     request.CacheName,
		Leaderboard:   request.LeaderboardName,
		ScoreRange:    scoreRange,
		Offset:        offset,
		LimitElements: count,
		Order:         leaderboardOrder,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return result.Elements, nil
}

func (client *leaderboardDataClient) GetRank(ctx context.Context, request *LeaderboardInternalGetRankRequest) ([]*pb.XRankedElement, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	leaderboardOrder := pb.XOrder_ASCENDING
	if request.Order != nil && *request.Order == DESCENDING {
		leaderboardOrder = pb.XOrder_DESCENDING
	}

	result, err := client.leaderboardClient.GetRank(ctx, &pb.XGetRankRequest{
		CacheName:   request.CacheName,
		Leaderboard: request.LeaderboardName,
		Ids:         request.Ids,
		Order:       leaderboardOrder,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return result.Elements, nil
}

func (client *leaderboardDataClient) Length(ctx context.Context, request *LeaderboardInternalLengthRequest) (uint32, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	result, err := client.leaderboardClient.GetLeaderboardLength(ctx, &pb.XGetLeaderboardLengthRequest{
		CacheName:   request.CacheName,
		Leaderboard: request.LeaderboardName,
	})
	if err != nil {
		return 0, momentoerrors.ConvertSvcErr(err)
	}
	return result.Count, nil
}

func (client *leaderboardDataClient) RemoveElements(ctx context.Context, request *LeaderboardInternalRemoveElementsRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	_, err := client.leaderboardClient.RemoveElements(ctx, &pb.XRemoveElementsRequest{
		CacheName:   request.CacheName,
		Leaderboard: request.LeaderboardName,
		Ids:         request.Ids,
	})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *leaderboardDataClient) Upsert(ctx context.Context, request *LeaderboardInternalUpsertRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	_, err := client.leaderboardClient.UpsertElements(ctx, &pb.XUpsertElementsRequest{
		CacheName:   request.CacheName,
		Leaderboard: request.LeaderboardName,
		Elements:    leaderboardUpsertElementToGrpc(request.Elements),
	})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
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
