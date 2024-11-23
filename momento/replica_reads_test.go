package momento_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("replica-reads", Label(CACHE_SERVICE_LABEL), func() {
	It(`should read the latest value after replication delay using balanced read concern`, func() {
		client, _ := sharedContext.GetClientPrereqsForType("withBalancedReadConcern")
		numTrials := 10
		delayBetweenTrials := 100 * time.Millisecond
		replicationDelay := 1 * time.Second

		trialFn := func(ctx context.Context, trialNumber int, wg *sync.WaitGroup, resultChan chan error) {
			defer wg.Done()

			// Start this trial with its own delay plus a random delay
			startDelay := time.Duration(trialNumber)*delayBetweenTrials + time.Duration(rand.Intn(20)-10)*time.Millisecond
			time.Sleep(startDelay)

			cacheKey := uuid.NewString()
			cacheValue := uuid.NewString()

			// Perform a set operation
			setResponse, err := client.Set(ctx, &momento.SetRequest{
				CacheName: sharedContext.CacheName,
				Key:       momento.String(cacheKey),
				Value:     momento.String(cacheValue),
			})

			Expect(setResponse).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			Expect(err).To(BeNil())

			// Wait for replication to complete
			time.Sleep(replicationDelay)

			// Verify that the value can be read
			getResponse, err := client.Get(ctx, &momento.GetRequest{
				CacheName: sharedContext.CacheName,
				Key:       momento.String(cacheKey),
			})
			Expect(getResponse).To(BeAssignableToTypeOf(&responses.GetHit{}))
			Expect(getResponse.(*responses.GetHit).ValueString()).To(Equal(cacheValue))
			Expect(err).To(BeNil())

			resultChan <- nil
		}

		// Run trials concurrently
		ctx := context.Background()
		var wg sync.WaitGroup
		resultChan := make(chan error, numTrials)

		for i := 0; i < numTrials; i++ {
			wg.Add(1)
			go trialFn(ctx, i, &wg, resultChan)
		}

		// Wait for all trials to complete
		wg.Wait()
		close(resultChan)

		// Collect and process results
		for err := range resultChan {
			Expect(err).To(BeNil(), fmt.Sprintf("Trial failed: %v", err))
		}
	})
})
