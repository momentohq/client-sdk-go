package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
)

var (
	publisherTopicName = os.Getenv("TEST_TOPIC_NAME")
)

func Publisher() {
	ctx := context.Background()
	credProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client, err := incubating.NewScsClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
	})
	if err != nil {
		panic(err)
	}
	err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "default",
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.AlreadyExistsError {
				panic(err)
			}
		}
	}
	fmt.Println(fmt.Sprintf("Publishing topic: %s", publisherTopicName))
	for {
		err = client.PublishTopic(ctx, &incubating.TopicPublishRequest{
			CacheName: "default",
			TopicName: publisherTopicName,
			Value:     strconv.Itoa(int(time.Now().UnixMilli())),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}