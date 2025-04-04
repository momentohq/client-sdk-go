package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"

	"github.com/google/uuid"
)

const (
	cacheName             = "my-test-cache"
	itemDefaultTTLSeconds = 60
)

type compressionMiddleware struct {
	middleware.Middleware
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func (mw *compressionMiddleware) GetRequestHandler(baseHandler middleware.RequestHandler) (middleware.RequestHandler, error) {
	return NewCompressionMiddlewareRequestHandler(baseHandler, mw.encoder, mw.decoder), nil
}

func NewCompressionMiddleware(props middleware.Props) middleware.Middleware {
	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
	decoder, _ := zstd.NewReader(nil)
	mw := middleware.NewMiddleware(props)
	return &compressionMiddleware{mw, encoder, decoder}
}

type compressionMiddlewareRequestHandler struct {
	middleware.RequestHandler
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewCompressionMiddlewareRequestHandler(rh middleware.RequestHandler, encoder *zstd.Encoder, decoder *zstd.Decoder) middleware.RequestHandler {
	return &compressionMiddlewareRequestHandler{rh, encoder, decoder}
}

func (rh *compressionMiddlewareRequestHandler) OnRequest(req interface{}) (interface{}, error) {
	// Compress on writes
	switch r := req.(type) {
	case *momento.SetRequest:
		rh.GetLogger().Info(fmt.Sprintf("(%s) Setting key: %s, value: %s", rh.GetId(), r.Key, r.Value))
		rawData := r.Value
		compressed := rh.encoder.EncodeAll([]byte(fmt.Sprintf("%v", rawData)), nil)
		rh.GetLogger().Info(
			fmt.Sprintf(
				"(%s) Compressed request %T: %d bytes -> %d bytes",
				rh.GetId(), req, len(fmt.Sprintf("%v", rawData)), len(compressed),
			),
		)
		return &momento.SetRequest{
			CacheName: r.CacheName,
			Key:       r.Key,
			Value:     momento.String(compressed),
		}, nil
	}
	return req, nil
}

func (rh *compressionMiddlewareRequestHandler) OnResponse(resp interface{}) (interface{}, error) {
	// Decompress on reads
	switch r := resp.(type) {
	case responses.GetHit:
		rawData := r.ValueByte()
		decompressed, err := rh.decoder.DecodeAll(rawData, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		rh.GetLogger().Info(
			fmt.Sprintf(
				"(%s) Decompressed response %T: %d bytes -> %d bytes",
				rh.GetId(), resp, len(rawData), len(decompressed),
			),
		)
		newGetResponse := responses.NewGetHit(decompressed)
		return newGetResponse, nil
	}
	return resp, nil
}

func doWork(ctx context.Context, client momento.CacheClient, index int) {
	// Sets key with default TTL and gets value with that key
	key := fmt.Sprintf("key-%d", index)
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err := client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
		Value:     momento.String(value),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err.Error())
	}

	switch r := resp.(type) {
	case *responses.GetHit:
		log.Printf("Lookup resulted in cache HIT. value=%s\n", r.ValueString())
	case *responses.GetMiss:
		log.Printf("Look up did not find a value key=%s", key)
	}
}

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	loggerFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)
	myConfig := config.LaptopLatest().WithMiddleware(
		[]middleware.Middleware{
			NewCompressionMiddleware(middleware.Props{Logger: loggerFactory.GetLogger("compression-middleware")}),
		},
	)

	// Initializes Momento
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		myConfig,
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}

	// Create Cache
	_, err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	// Make some requests
	for i := 0; i < 5; i++ {
		doWork(ctx, client, i)
	}

	// Permanently delete the cache
	if _, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName}); err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}
