package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/utils"
)

type PreviewLeaderboardClient interface {
	Leaderboard(ctx context.Context, request *LeaderboardRequest) (Leaderboard, error)
	Close()
}

type previewLeaderboardClient struct {
	credentialProvider    auth.CredentialProvider
	leaderboardDataClient *leaderboardDataClient
	log                   logger.MomentoLogger
}

func NewPreviewLeaderboardClient(leaderboardConfiguration config.LeaderboardConfiguration, credentialProvider auth.CredentialProvider) (PreviewLeaderboardClient, error) {
	dataClient, err := newLeaderboardDataClient(&models.LeaderboardClientRequest{
		CredentialProvider: credentialProvider,
		Configuration:      leaderboardConfiguration,
	})
	if err != nil {
		return nil, err
	}

	client := &previewLeaderboardClient{
		credentialProvider:    credentialProvider,
		leaderboardDataClient: dataClient,
		log:                   leaderboardConfiguration.GetLoggerFactory().GetLogger("topic-client"),
	}
	return client, nil
}

func (c previewLeaderboardClient) Leaderboard(ctx context.Context, request *LeaderboardRequest) (Leaderboard, error) {
	if err := utils.ValidateName(request.CacheName, "cache name"); err != nil {
		return nil, err
	}
	if err := utils.ValidateName(request.LeaderboardName, "leaderboard name"); err != nil {
		return nil, err
	}
	newLeaderboard := &leaderboard{
		cacheName:             request.CacheName,
		leaderboardName:       request.LeaderboardName,
		leaderboardDataClient: c.leaderboardDataClient,
	}
	return newLeaderboard, nil
}

func (c previewLeaderboardClient) Close() {
	c.leaderboardDataClient.close()
}
