package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/config/retry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// AddUnaryRetryInterceptor returns a unary interceptor that will retry the request based on the retry strategy.
func AddUnaryRetryInterceptor(s retry.Strategy, onRequest func(context.Context, string)) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		attempt := 1
		for {
			fmt.Printf("unary interceptor attempt %d for %s\n", attempt, method)
			// This is currently used for testing purposes only by the RetryMetricsMiddleware.
			if onRequest != nil {
				onRequest(ctx, method)
			}
			// Execute api call
			lastErr := invoker(ctx, method, req, reply, cc, opts...)
			if lastErr == nil {
				fmt.Printf("got REPLY -----> %#v\n", reply)
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

// AddStreamRetryInterceptor returns a stream interceptor that will retry the request based on the retry strategy.
func AddStreamRetryInterceptor(s retry.Strategy, onRequest func(context.Context, string)) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		attempt := 1
		for {
			fmt.Printf("attempt %d\n", attempt)
			// This is currently used for testing purposes only by the RetryMetricsMiddleware.
			if onRequest != nil {
				onRequest(ctx, method)
			}
			// Execute api call
			stream, lastErr := streamer(ctx, desc, cc, method, opts...)
			fmt.Printf("lastErr: %v\n", lastErr)
			if lastErr == nil {
				//fmt.Printf("no subscribe error, returning %v", stream)
				// Success no error returned stop interceptor
				return stream, nil
			}

			// Check retry eligibility based off last error received
			retryBackoffTime := s.DetermineWhenToRetry(retry.StrategyProps{
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
				time.Sleep(time.Duration(*retryBackoffTime) * time.Millisecond)
			}
			attempt++
		}
	}
}
