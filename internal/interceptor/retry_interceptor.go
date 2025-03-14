package interceptor

import (
	"context"
	"github.com/momentohq/client-sdk-go/config/retry"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func AddUnaryRetryInterceptor(s retry.Strategy) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		attempt := 1
		for {

			// Execute api call
			lastErr := invoker(ctx, method, req, reply, cc, opts...)
			if lastErr == nil {
				// Success no error returned stop interceptor
				return nil
			}

			// Check retry eligibility based off last error received
			retryBackoffTime := s.DetermineWhenToRetry(retry.StrategyProps{
				GrpcStatusCode: status.Code(lastErr),
				GrpcMethod:     method,
				AttemptNumber:  attempt,
			})

			if retryBackoffTime == nil {
				// If nil backoff time don't retry just return last error received
				return lastErr
			}

			// Sleep for recommended time interval and increment attempts before trying again
			if *retryBackoffTime > 0 {
				time.Sleep(time.Duration(*retryBackoffTime) * time.Millisecond)
			}
			attempt++
		}
	}
}
