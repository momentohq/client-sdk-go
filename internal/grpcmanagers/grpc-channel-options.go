package grpcmanagers

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/momentohq/client-sdk-go/config"
)

func GrpcChannelOptionsFromGrpcConfig(grpcConfig config.GrpcConfiguration) []grpc.DialOption {
	if grpcConfig == nil {
		return []grpc.DialOption{}
	}

	return []grpc.DialOption{
		// max call sizes are CallOptions
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(grpcConfig.GetMaxReceiveMessageLength()),
			grpc.MaxCallSendMsgSize(grpcConfig.GetMaxSendMessageLength()),
		),

		// keepalive params are DialOptions
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                grpcConfig.GetKeepAliveTime(),
				Timeout:             grpcConfig.GetKeepAliveTimeout(),
				PermitWithoutStream: grpcConfig.GetKeepAlivePermitWithoutCalls(),
			},
		),
	}
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
