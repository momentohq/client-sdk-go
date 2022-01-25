package grpcmanager

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type DataGrpcManager struct {
	Conn *grpc.ClientConn
}

func NewDataGrpcManager(authToken string, endPoint string) (DataGrpcManager, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(addHeadersInterceptorData(authToken)))
	return DataGrpcManager{Conn: conn}, err
}

func (cm *DataGrpcManager) Close() error {
	return cm.Conn.Close()
}

func addHeadersInterceptorData(authToken string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authToken)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
