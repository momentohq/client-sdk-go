package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/momento"
)

const (
	DefaultClient             = "defaultClient"
	WithDefaultCache          = "withDefaultCache"
	WithConsistentReadConcern = "withConsistentReadConcern"
	WithBalancedReadConcern   = "withBalancedReadConcern"
)

var consistentReads bool

func init() {
	consistentReads = os.Getenv("CONSISTENT_READS") != ""
	log.Printf("consistent reads set to %t for integration tests", consistentReads)
}

type SharedContext struct {
	Client                          momento.CacheClient
	ClientWithDefaultCacheName      momento.CacheClient
	ClientWithConsistentReadConcern momento.CacheClient
	ClientWithBalancedReadConcern   momento.CacheClient
	DefaultCacheName                string
	TopicClient                     momento.TopicClient
	CacheName                       string
	Ctx                             context.Context
	DefaultTtl                      time.Duration
	Configuration                   config.Configuration
	TopicConfiguration              config.TopicsConfiguration
	CredentialProvider              auth.CredentialProvider
	AuthClient                      momento.AuthClient
	AuthConfiguration               config.AuthConfiguration
	LeaderboardClient               momento.PreviewLeaderboardClient
	LeaderboardConfiguration        config.LeaderboardConfiguration
}

type SharedContextProps struct {
	IsMomentoLocal bool
}

func NewSharedContext(props SharedContextProps) SharedContext {
	shared := SharedContext{}
	shared.Ctx = context.Background()
	var credentialProvider auth.CredentialProvider
	var err error
	if props.IsMomentoLocal {
		port := os.Getenv("MOMENTO_LOCAL_PORT")
		if port == "" {
			port = "8080"
		}
		thePort, err := strconv.ParseUint(port, 10, 32)
		if err != nil {
			panic(fmt.Sprintf("invalid port %v", port))
		}
		credentialProvider, err = auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{
			Hostname: "",
			Port:     uint(thePort),
		})
		if err != nil {
			panic(err)
		}
	} else {
		//lint:ignore SA1019 // Still supporting FromEnvironmentVariable for backwards compatibility
		credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
		if err != nil {
			panic(err)
		}
	}
	shared.CredentialProvider = credentialProvider
	shared.Configuration = config.LaptopLatestWithLogger(logger.NewNoopMomentoLoggerFactory()).WithClientTimeout(15 * time.Second)
	shared.TopicConfiguration = config.TopicsDefaultWithLogger(logger.NewNoopMomentoLoggerFactory())
	shared.AuthConfiguration = config.AuthDefaultWithLogger(logger.NewNoopMomentoLoggerFactory())
	shared.LeaderboardConfiguration = config.LeaderboardDefaultWithLogger(logger.NewNoopMomentoLoggerFactory())
	shared.DefaultTtl = 3 * time.Second

	var clientConfig config.Configuration
	if consistentReads {
		clientConfig = shared.Configuration.WithReadConcern(config.CONSISTENT)
	} else {
		clientConfig = shared.Configuration
	}

	client, err := momento.NewCacheClient(clientConfig, shared.CredentialProvider, shared.DefaultTtl)

	if err != nil {
		panic(err)
	}

	defaultCacheName := fmt.Sprintf("golang-default-%s", uuid.NewString())
	clientDefaultCacheName, err := momento.NewCacheClientWithDefaultCache(
		clientConfig, shared.CredentialProvider, shared.DefaultTtl, defaultCacheName,
	)
	if err != nil {
		panic(err)
	}

	consistentReadConcernClient, err := momento.NewCacheClientWithDefaultCache(
		shared.Configuration.WithReadConcern(config.CONSISTENT), shared.CredentialProvider, shared.DefaultTtl, defaultCacheName,
	)
	if err != nil {
		panic(err)
	}

	balancedReadConcernClient, err := momento.NewCacheClientWithDefaultCache(
		shared.Configuration.WithReadConcern(config.BALANCED), shared.CredentialProvider, shared.DefaultTtl, defaultCacheName,
	)
	if err != nil {
		panic(err)
	}

	topicClient, err := momento.NewTopicClient(shared.TopicConfiguration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	authClient, err := momento.NewAuthClient(shared.AuthConfiguration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	leaderboardClient, err := momento.NewPreviewLeaderboardClient(shared.LeaderboardConfiguration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	shared.Client = client
	shared.ClientWithDefaultCacheName = clientDefaultCacheName
	shared.ClientWithConsistentReadConcern = consistentReadConcernClient
	shared.ClientWithBalancedReadConcern = balancedReadConcernClient
	shared.DefaultCacheName = defaultCacheName
	shared.TopicClient = topicClient
	shared.AuthClient = authClient
	shared.LeaderboardClient = leaderboardClient

	shared.CacheName = fmt.Sprintf("golang-%s", uuid.NewString())

	return shared
}

func (shared SharedContext) GetClientPrereqsForType(clientType string) (momento.CacheClient, string) {
	var client momento.CacheClient
	var cacheName string
	if clientType == WithDefaultCache {
		client = shared.ClientWithDefaultCacheName
		cacheName = ""
	} else if clientType == DefaultClient {
		client = shared.Client
		cacheName = shared.CacheName
	} else if clientType == WithConsistentReadConcern {
		client = shared.ClientWithConsistentReadConcern
		cacheName = ""
	} else if clientType == WithBalancedReadConcern {
		client = shared.ClientWithBalancedReadConcern
	} else {
		panic("invalid client type")
	}
	return client, cacheName
}

func (shared SharedContext) CreateDefaultCaches() {
	_, err := shared.Client.CreateCache(shared.Ctx, &momento.CreateCacheRequest{CacheName: shared.CacheName})
	if err != nil {
		panic(err)
	}
	_, err = shared.Client.CreateCache(shared.Ctx, &momento.CreateCacheRequest{CacheName: shared.DefaultCacheName})
	if err != nil {
		panic(err)
	}
}

func (shared SharedContext) Close() {
	// close topic client before deleting the cache in which it is subscribed
	shared.TopicClient.Close()

	_, err := shared.Client.DeleteCache(shared.Ctx, &momento.DeleteCacheRequest{CacheName: shared.CacheName})
	if err != nil {
		panic(err)
	}
	_, err = shared.Client.DeleteCache(shared.Ctx, &momento.DeleteCacheRequest{CacheName: shared.DefaultCacheName})
	if err != nil {
		panic(err)
	}

	shared.Client.Close()
	shared.AuthClient.Close()
	shared.LeaderboardClient.Close()
}
