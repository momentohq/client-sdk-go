package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

// Defines the set of operations that can be performed on a leaderboard object.
type Leaderboard interface {
	// Deletes the leaderboard, i.e. removes all elements from the leaderboard.
	Delete(ctx context.Context) (responses.LeaderboardDeleteResponse, error)

	// Fetches elements that fall within the specified min and max ranks.
	FetchByRank(ctx context.Context, request LeaderboardFetchByRankRequest) (responses.LeaderboardFetchResponse, error)

	// Fetches elements that fall within the specified min and max scores. Elements with the same score will be
	// returned in alphanumerical order based on their ID (e.g. IDs of elements with the same score would be
	// returned in the order [1, 10, 123, 2, 234, ...] rather than [1, 2, 10, 123, 234, ...]).
	FetchByScore(ctx context.Context, request LeaderboardFetchByScoreRequest) (responses.LeaderboardFetchResponse, error)

	// Fetches elements (including rank, score, and ID) given a list of element IDs.
	GetRank(ctx context.Context, request LeaderboardGetRankRequest) (responses.LeaderboardFetchResponse, error)

	// Gets the number of entries in the leaderboard.
	Length(ctx context.Context) (responses.LeaderboardLengthResponse, error)

	// Removes elements with the specified IDs from the leaderboard.
	RemoveElements(ctx context.Context, request LeaderboardRemoveElementsRequest) (responses.LeaderboardRemoveElementsResponse, error)

	// Inserts elements if they do not already exist in the leaderboard. Updates elements if they do already exist
	// in the leaderboard. There are no partial failures; an upsert call will either succeed or fail.
	Upsert(ctx context.Context, request LeaderboardUpsertRequest) (responses.LeaderboardUpsertResponse, error)
}

type leaderboard struct {
	cacheName             string
	leaderboardName       string
	leaderboardDataClient *leaderboardDataClient
}

// Deletes the leaderboard, i.e. removes all elements from the leaderboard.
func (l *leaderboard) Delete(ctx context.Context) (responses.LeaderboardDeleteResponse, error) {
	r := &LeaderboardInternalDeleteRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
	}
	if err := l.leaderboardDataClient.delete(ctx, r); err != nil {
		return nil, err
	}
	return &responses.LeaderboardDeleteSuccess{}, nil
}

// Fetches elements that fall within the specified min and max ranks.
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
	elements, err := l.leaderboardDataClient.fetchByRank(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardFetchSuccess(leaderboardFetchGrpcElementToModel(elements)), nil
}

// Fetches elements that fall within the specified min and max scores. Elements with the same score will be
// returned in alphanumerical order based on their ID (e.g. IDs of elements with the same score would be
// returned in the order [1, 10, 123, 2, 234, ...] rather than [1, 2, 10, 123, 234, ...]).
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
	elements, err := l.leaderboardDataClient.fetchByScore(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardFetchSuccess(leaderboardFetchGrpcElementToModel(elements)), nil
}

// Fetches elements (including rank, score, and ID) given a list of element IDs.
func (l *leaderboard) GetRank(ctx context.Context, request LeaderboardGetRankRequest) (responses.LeaderboardFetchResponse, error) {
	r := &LeaderboardInternalGetRankRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		Ids:             request.Ids,
		Order:           request.Order,
	}
	elements, err := l.leaderboardDataClient.getRank(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardFetchSuccess(leaderboardFetchGrpcElementToModel(elements)), nil
}

// Gets the number of entries in the leaderboard.
func (l *leaderboard) Length(ctx context.Context) (responses.LeaderboardLengthResponse, error) {
	r := &LeaderboardInternalLengthRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
	}
	length, err := l.leaderboardDataClient.length(ctx, r)
	if err != nil {
		return nil, err
	}
	return responses.NewLeaderboardLengthSuccess(length), nil
}

// Removes elements with the specified IDs from the leaderboard.
func (l *leaderboard) RemoveElements(ctx context.Context, request LeaderboardRemoveElementsRequest) (responses.LeaderboardRemoveElementsResponse, error) {
	if len(request.Ids) == 0 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "List of elements to remove cannot be empty", nil)
	}
	r := &LeaderboardInternalRemoveElementsRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		Ids:             request.Ids,
	}
	if err := l.leaderboardDataClient.removeElements(ctx, r); err != nil {
		return nil, err
	}
	return &responses.LeaderboardRemoveElementsSuccess{}, nil
}

// Inserts elements if they do not already exist in the leaderboard. Updates elements if they do already exist
// in the leaderboard. There are no partial failures; an upsert call will either succeed or fail.
func (l *leaderboard) Upsert(ctx context.Context, request LeaderboardUpsertRequest) (responses.LeaderboardUpsertResponse, error) {
	if len(request.Elements) == 0 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "List of elements to upsert cannot be empty", nil)
	}
	r := &LeaderboardInternalUpsertRequest{
		CacheName:       l.cacheName,
		LeaderboardName: l.leaderboardName,
		Elements:        request.Elements,
	}
	if err := l.leaderboardDataClient.upsert(ctx, r); err != nil {
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
