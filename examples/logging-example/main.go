package main

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/go-example/logging-example/momento_logrus"
	logrus "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println(`
This example illustrates how to construct a Momento client that is configured to use logrus for
logging output.  You can use a similar approach to integrate with the logging framework of your
choice, so that Momento's logs will be written to the same destinations as the rest of your
application's logs.`)
	fmt.Println("")

	creds, err := auth.FromEnvironmentVariable("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatestWithLogger(momento_logrus.NewLogrusMomentoLoggerFactory()),
		creds,
		60*time.Second,
		30*time.Second,
	)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	fmt.Println("Creating and deleting a cache with logrus level set to FATAL; you should not see any momento log messages.")
	fmt.Println("")
	logrus.SetLevel(logrus.FatalLevel)
	createAndDeleteCache(ctx, client)
	fmt.Println("Cache created and deleted!")
	fmt.Println("")
	fmt.Println("Adjusting logrus log level to INFO")
	logrus.SetLevel(logrus.InfoLevel)
	fmt.Println("Creating and deleting a cache with logrus level set to INFO; this time you SHOULD see some momento log messages.")
	fmt.Println("")
	createAndDeleteCache(ctx, client)

	fmt.Println(`
Done!  See momento_logrus/logrus.go for the details of how we connected the Momento logging to logrus. If you need
a similar integration for another logging framework such as zerolog, zap, etc. please feel free to reach out to us!
`)
}

func createAndDeleteCache(ctx context.Context, client momento.CacheClient) {
	cacheName := "logging-example-cache"
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	_, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}
}
