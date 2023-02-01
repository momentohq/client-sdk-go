package main

import (
	"context"
	"errors"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/prozz/aws-embedded-metrics-golang/emf"
	"time"
)

func createCacheIfNotExist(ctx context.Context, client momento.ScsClient, cacheName string) {
	err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.AlreadyExistsError {
				panic(err)
			}
		}
	}
}

func sendSubscriberCountToCw() {
	go func() {
		for {
			emf.New(emf.WithLogGroup("pubsub")).MetricAs("SubscriberCount", 1, emf.Count).DimensionSet(emf.NewDimension("subscriber", "receiving-test")).Log()
			time.Sleep(time.Minute)
		}
	}()
}

func sendLatencyToCw(latency int) {
	emf.New(emf.WithLogGroup("pubsub")).MetricAs("ReceivingMessageLatency", latency, emf.Milliseconds).DimensionSet(emf.NewDimension("subscriber", "receiving-test")).Log()
}
