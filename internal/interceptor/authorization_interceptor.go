package interceptor

import (
	"context"
	"github.com/momentohq/client-sdk-go/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AddHeadersInterceptor(authToken string, readConcern config.ReadConcern) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "read_concern", string(readConcern), "authorization", authToken), method, req, reply, cc, opts...)
	}
}

func AddStreamHeaderInterceptor(authToken string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "authorization", authToken), desc, cc, method)
	}
}
