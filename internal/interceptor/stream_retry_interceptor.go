package interceptor

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// streamInterceptor is an example stream interceptor.
func StreamInterceptor(retryStrat retry.Strategy) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		attempt := 1
		for {
			s, lastErr := streamer(ctx, desc, cc, method, opts...)
			if lastErr == nil {
				// Success no error returned stop interceptor
				return s, nil
			}

			// Check retry eligibility based off last error received
			retryBackoffTime := retryStrat.DetermineWhenToRetry(retry.StrategyProps{
				GrpcStatusCode: status.Code(lastErr),
				GrpcMethod:     method,
				AttemptNumber:  attempt,
			})

			if retryBackoffTime == nil {
				// If nil backoff time don't retry just return last error received
				return nil, lastErr
			}

			// Sleep for recommended time interval and increment attempts before trying again
			if *retryBackoffTime > 0 {
				time.Sleep(time.Duration(*retryBackoffTime) * time.Second)
			}
			attempt++
		}
	}
}
