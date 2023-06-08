package interceptor

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// StreamInterceptor is an interceptor that handles retries for streaming apis
func StreamRetryInterceptor(retryStrat retry.Strategy) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		attempt := 1
		for {
			s, lastErr := streamer(ctx, desc, cc, method, opts...)
			if lastErr == nil {
				// Successfully connected the stream. If we wanted to get fancy here, we could return 
				// a struct that implements the Stream interface, and wraps the SendMsg and ReceiveMsg function
				// calls. Since we only care about retrying on subscribe, we just return the stream here
				// ex. https://github.com/grpc/grpc-go/blob/7a7caf363d9b33bfa5ddac83e7dab744f695fb5b/examples/features/interceptor/server/main.go#L106-L141
				return s, nil
			}

			// Unable to connect to the stream, checking retry eligibility based off last error received
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
