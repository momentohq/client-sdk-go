package interceptor

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var FirstTimeHeadersSent = false
var Version = "1.24.0" // x-release-please-version

func AddAgentHeaderInterceptor(clientType string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "Agent", fmt.Sprintf("golang:%s:%s", clientType, Version)), method, req, reply, cc, opts...)
	}
}

func AddRuntimeVersionHeaderInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "Runtime-Version", fmt.Sprintf("golang:%s", runtime.Version())), method, req, reply, cc, opts...)
	}
}

func AddStreamAgentHeaderInterceptor(clientType string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "Agent", fmt.Sprintf("golang:%s:%s", clientType, Version)), desc, cc, method)
	}
}

func AddStreamRuntimeVersionHeaderInterceptor() func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "Runtime-Version", fmt.Sprintf("golang:%s", runtime.Version())), desc, cc, method)
	}
}
