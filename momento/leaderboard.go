package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type Leaderboard interface {
	Delete(ctx context.Context) (responses.LeaderboardDeleteResponse, error)
	FetchByRank(ctx context.Context, request LeaderboardFetchByRankRequest) (responses.LeaderboardFetchResponse, error)
	FetchByScore(ctx context.Context, request LeaderboardFetchByScoreRequest) (responses.LeaderboardFetchResponse, error)
	GetRank(ctx context.Context, request LeaderboardGetRankRequest) (responses.LeaderboardFetchResponse, error)
	Length(ctx context.Context) (responses.LeaderboardLengthResponse, error)
	RemoveElements(ctx context.Context, request LeaderboardRemoveElementsRequest) (responses.LeaderboardRemoveElementsResponse, error)
	Upsert(ctx context.Context, request LeaderboardUpsertRequest) (responses.LeaderboardUpsertResponse, error)
}

type leaderboard struct {
	cacheName             string
	leaderboardName       string
	leaderboardDataClient *leaderboardDataClient
}

// Delete implements Leaderboard.
func (l *leaderboard) Delete(ctx context.Context) (responses.LeaderboardDeleteResponse, error) {
	r := &LeaderboardInternalDeleteRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
	}
	if err := l.leaderboardDataClient.Delete(ctx, r); err != nil {
		return nil, err
	}
	return &responses.LeaderboardDeleteSuccess{}, nil
}

// FetchByRank implements Leaderboard.
func (l *leaderboard) FetchByRank(ctx context.Context, request LeaderboardFetchByRankRequest) (responses.LeaderboardFetchResponse, error) {
	if request.StartRank >= request.EndRank {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "start rank must be less than end rank", nil)
	}
	r := &LeaderboardInternalFetchByRankRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		StartRank:       request.StartRank,
		EndRank:         request.EndRank,
		Order:           request.Order,
	}
	elements, err := l.leaderboardDataClient.FetchByRank(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardFetchSuccess(leaderboardFetchGrpcElementToModel(elements)), nil
}

// FetchByScore implements Leaderboard.
func (l *leaderboard) FetchByScore(ctx context.Context, request LeaderboardFetchByScoreRequest) (responses.LeaderboardFetchResponse, error) {
	if request.MinScore != nil && request.MaxScore != nil && *request.MinScore >= *request.MaxScore {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "min score must be less than max score", nil)
	}
	r := &LeaderboardInternalFetchByScoreRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		MinScore:        request.MinScore,
		MaxScore:        request.MaxScore,
		Offset:          request.Offset,
		Count:           request.Count,
		Order:           request.Order,
	}
	elements, err := l.leaderboardDataClient.FetchByScore(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardFetchSuccess(leaderboardFetchGrpcElementToModel(elements)), nil
}

// GetRank implements Leaderboard.
func (l *leaderboard) GetRank(ctx context.Context, request LeaderboardGetRankRequest) (responses.LeaderboardFetchResponse, error) {
	r := &LeaderboardInternalGetRankRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		Ids:             request.Ids,
		Order:           request.Order,
	}
	elements, err := l.leaderboardDataClient.GetRank(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardFetchSuccess(leaderboardFetchGrpcElementToModel(elements)), nil
}

// Length implements Leaderboard.
func (l *leaderboard) Length(ctx context.Context) (responses.LeaderboardLengthResponse, error) {
	r := &LeaderboardInternalLengthRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
	}
	length, err := l.leaderboardDataClient.Length(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardLengthSuccess(length), nil
}

// RemoveElements implements Leaderboard.
func (l *leaderboard) RemoveElements(ctx context.Context, request LeaderboardRemoveElementsRequest) (responses.LeaderboardRemoveElementsResponse, error) {
	if len(request.Ids) == 0 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "List of elements to remove cannot be empty", nil)
	}
	r := &LeaderboardInternalRemoveElementsRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		Ids:             request.Ids,
	}
	if err := l.leaderboardDataClient.RemoveElements(ctx, r); err != nil {
		return nil, err
	}
	return &responses.LeaderboardRemoveElementsSuccess{}, nil
}

// Upsert implements Leaderboard.
func (l *leaderboard) Upsert(ctx context.Context, request LeaderboardUpsertRequest) (responses.LeaderboardUpsertResponse, error) {
	if len(request.Elements) == 0 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "List of elements to upsert cannot be empty", nil)
	}
	r := &LeaderboardInternalUpsertRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		Elements:        request.Elements,
	}
	if err := l.leaderboardDataClient.Upsert(ctx, r); err != nil {
		return nil, err
	}
	return &responses.LeaderboardUpsertSuccess{}, nil
}

func leaderboardFetchGrpcElementToModel(grpcRankedElements []*pb.XRankedElement) []responses.LeaderboardElement {
	var returnList []responses.LeaderboardElement
	for _, element := range grpcRankedElements {
		returnList = append(returnList, responses.LeaderboardElement{
			Id:    element.Id,
			Rank:  element.Rank,
			Score: element.Score,
		})
	}
	return returnList
}
