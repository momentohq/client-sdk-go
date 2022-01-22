package grpcmanager

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type ControlGrpcManager struct {
	Conn	*grpc.ClientConn
}

func NewControlGrpcManager(authToken string, endPoint string) (ControlGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(addHeadersInterceptorControl(authToken)))
	return ControlGrpcManager{Conn: conn}, err
}

func (cm *ControlGrpcManager) Close() error {
	return cm.Conn.Close()
}

func addHeadersInterceptorControl(authToken string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func (ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authToken)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
