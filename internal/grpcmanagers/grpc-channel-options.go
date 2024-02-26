package grpcmanagers

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/momentohq/client-sdk-go/config"
)

// The default value for max_send_message_length is 4mb.  We need to increase this to 5mb in order to
// support cases where users have requested a limit increase up to our maximum item size of 5mb.
const DEFAULT_MAX_MESSAGE_SIZE = 5_243_000

func GrpcChannelOptionsFromGrpcConfig(grpcConfig config.GrpcConfiguration) []grpc.DialOption {
	if grpcConfig == nil {
		return []grpc.DialOption{
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(DEFAULT_MAX_MESSAGE_SIZE),
				grpc.MaxCallSendMsgSize(DEFAULT_MAX_MESSAGE_SIZE),
			),
		}
	}

	options := make([]grpc.DialOption, 0, 2)

	max_send_length := DEFAULT_MAX_MESSAGE_SIZE
	if *grpcConfig.GetMaxSendMessageLength() > 0 {
		max_send_length = *grpcConfig.GetMaxSendMessageLength()
	}

	max_receive_length := DEFAULT_MAX_MESSAGE_SIZE
	if *grpcConfig.GetMaxReceiveMessageLength() > 0 {
		max_receive_length = *grpcConfig.GetMaxReceiveMessageLength()
	}

	options = append(options, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(max_receive_length),
		grpc.MaxCallSendMsgSize(max_send_length),
	))

	// If keepAlivePermitWithoutCalls is not defined in the config, a default value
	// of false is used. We will only set the keepalive settings if keepAlivePermitWithoutCalls
	// is set to true.
	if *grpcConfig.GetKeepAlivePermitWithoutCalls() {
		options = append(options, grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                *grpcConfig.GetKeepAliveTime(),
				Timeout:             *grpcConfig.GetKeepAliveTimeout(),
				PermitWithoutStream: *grpcConfig.GetKeepAlivePermitWithoutCalls(),
			},
		))
	}

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
