package batchutils_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/batchutils"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	. "github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Batch operations", func() {

	var (
		ctx       context.Context
		client    CacheClient
		cacheName string
		keys      []Value
	)

	BeforeEach(func() {
		ctx = context.Background()
		cacheName = fmt.Sprintf("golang-%s", uuid.NewString())
		credentialProvider, err := auth.FromEnvironmentVariable("TEST_AUTH_TOKEN")
		if err != nil {
			panic(err)
		}
		client, err = NewCacheClient(
			config.LaptopLatestWithLogger(logger.NewNoopMomentoLoggerFactory()),
			credentialProvider,
			time.Second*60,
		)
		if err != nil {
			panic(err)
		}

		_, err = client.CreateCache(ctx, &CreateCacheRequest{CacheName: cacheName})
		if err != nil {
			panic(err)
		}

		for i := 0; i < 50; i++ {
			key := String(fmt.Sprintf("key%d", i))
			keys = append(keys, key)
			_, err := client.Set(ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String(fmt.Sprintf("val%d", i)),
			})
			if err != nil {
				panic(err)
			}
		}
	})

	AfterEach(func() {
		_, err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: cacheName})
		if err != nil {
			panic(err)
		}
	})

	It("batch deletes", func() {
		errors := batchutils.BatchDelete(ctx, &batchutils.BatchDeleteRequest{
			Client:    client,
			CacheName: cacheName,
			Keys:      keys[5:21],
		})
		Expect(len(errors)).To(Equal(0))
		for i := 0; i < 50; i++ {
			resp, err := client.Get(ctx, &GetRequest{
				CacheName: cacheName,
				Key:       keys[i],
			})
			Expect(err).To(BeNil())
			switch resp.(type) {
			case *responses.GetHit:
				if i >= 5 && i <= 20 {
					Fail("got a hit for #%d that should be a miss", i)
				}
			case *responses.GetMiss:
				if !(i >= 5 && i <= 20) {
					Fail("got a miss for #%d that should be a hit", i)
				}
			}
		}
	})
})
