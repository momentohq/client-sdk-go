package interceptor

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/config/retry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// AddUnaryRetryInterceptor returns a unary interceptor that will retry the request based on the retry strategy.
func AddUnaryRetryInterceptor(s retry.Strategy, onRequest func(context.Context, string)) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		attempt := 1
		var overallDeadline time.Time
		deadline, ok := ctx.Deadline()
		if ok {
			overallDeadline = deadline
		} else {
			overallDeadline = time.Now().Add(5 * time.Second)
		}

		for {
			// This is currently used for testing purposes only by the RetryMetricsMiddleware.
			if onRequest != nil {
				onRequest(ctx, method)
			}

			// If the FixedTimeoutRetryStrategy is used, overwrite the deadline using
			// the configured retry timeout. Otherwise, use the given context.
			retryCtx := ctx
			if s, ok := s.(retry.FixedTimeoutRetryStrategy); ok {
				retryDeadline := s.CalculateRetryDeadline(overallDeadline)
				ctxWithRetryDeadline, cancel := context.WithDeadline(ctx, retryDeadline)
				defer cancel()
				retryCtx = ctxWithRetryDeadline
			}

			// Execute api call
			lastErr := invoker(retryCtx, method, req, reply, cc, opts...)
			if lastErr == nil {
				// Success no error returned stop interceptor
				return nil
			}

			if s == nil {
				// No retry strategy is configured so return the error
				return lastErr
			}

			// Check retry eligibility based off last error received
			retryBackoffTime := s.DetermineWhenToRetry(retry.StrategyProps{
				GrpcStatusCode:  status.Code(lastErr),
				GrpcMethod:      method,
				AttemptNumber:   attempt,
				OverallDeadline: overallDeadline,
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

type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// AddStreamRetryInterceptor returns a stream interceptor that will wrap the stream for inspection.
// This is currently unused but I want to keep it here for reference.
func AddStreamInterceptor() func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		return newWrappedStream(s), nil
	}
}
