package grpcmanagers

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/momentohq/client-sdk-go/config"
)

// The default value for max_send_message_length is 4mb.  We need to increase this to 5mb in order to
// support cases where users have requested a limit increase up to our maximum item size of 5mb.
const DEFAULT_MAX_MESSAGE_SIZE = 5_243_000

func GrpcChannelOptionsFromGrpcConfig(grpcConfig config.GrpcConfiguration) []grpc.DialOption {
	// If no grpc config is provided, we default to 5mb message sizes and keepalive turned on.
	if grpcConfig == nil {
		return []grpc.DialOption{
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(DEFAULT_MAX_MESSAGE_SIZE),
				grpc.MaxCallSendMsgSize(DEFAULT_MAX_MESSAGE_SIZE),
			),
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					PermitWithoutStream: true,
					Time:                5000 * time.Millisecond,
					Timeout:             1000 * time.Millisecond,
				},
			),
		}
	}

	// Otherwise construct the options from the provided config.

	options := make([]grpc.DialOption, 0, 2)

	max_send_length := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfig.GetMaxSendMessageLength() > 0 {
		max_send_length = grpcConfig.GetMaxSendMessageLength()
	}

	max_receive_length := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfig.GetMaxReceiveMessageLength() > 0 {
		max_receive_length = grpcConfig.GetMaxReceiveMessageLength()
	}

	options = append(options, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(max_receive_length),
		grpc.MaxCallSendMsgSize(max_send_length),
	))

	keepaliveOptions := keepalive.ClientParameters{}

	if grpcConfig.GetKeepAlivePermitWithoutCalls() {
		keepaliveOptions.PermitWithoutStream = true
	}

	if grpcConfig.GetKeepAliveTime() > 0 {
		keepaliveOptions.Time = grpcConfig.GetKeepAliveTime()
	}

	if grpcConfig.GetKeepAliveTimeout() > 0 {
		keepaliveOptions.Timeout = grpcConfig.GetKeepAliveTimeout()
	}

	options = append(options, grpc.WithKeepaliveParams(keepaliveOptions))

	return options
}

func TransportCredentialsChannelOption() grpc.DialOption {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(config))
}

func AllDialOptions(grpcConfig config.GrpcConfiguration, options ...grpc.DialOption) []grpc.DialOption {
	options = append(options, TransportCredentialsChannelOption())
	options = append(options, GrpcChannelOptionsFromGrpcConfig(grpcConfig)...)
	return options
}
