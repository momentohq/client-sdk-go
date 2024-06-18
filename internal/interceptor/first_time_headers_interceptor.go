package interceptor

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var FirstTimeHeadersSent = false

func AddAgentHeaderInterceptor(clientType string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	version := "1.22.0"
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "Agent", fmt.Sprintf("golang:%s:%s", clientType, version)), method, req, reply, cc, opts...)
	}
}

func AddRuntimeVersionHeaderInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "Runtime-Version", fmt.Sprintf("golang:%s", runtime.Version())), method, req, reply, cc, opts...)
	}
}

func AddStreamAgentHeaderInterceptor(clientType string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	version := "1.22.0"
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "Agent", fmt.Sprintf("golang:%s:%s", clientType, version)), desc, cc, method)
	}
}

func AddStreamRuntimeVersionHeaderInterceptor() func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "Runtime-Version", fmt.Sprintf("golang:%s", runtime.Version())), desc, cc, method)
	}
}
