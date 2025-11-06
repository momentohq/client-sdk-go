package grpcmanagers

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
)

func GrpcChannelOptionsFromGrpcConfig(grpcConfig config.IGrpcConfiguration) []grpc.DialOption {
	// Default to 5mb message sizes and keepalives turned on (defaults are set in NewStaticGrpcConfiguration)
	return []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(grpcConfig.GetMaxReceiveMessageLength()),
			grpc.MaxCallSendMsgSize(grpcConfig.GetMaxSendMessageLength()),
		),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				PermitWithoutStream: grpcConfig.GetKeepAlivePermitWithoutCalls(),
				Time:                grpcConfig.GetKeepAliveTime(),
				Timeout:             grpcConfig.GetKeepAliveTimeout(),
			},
		),
		grpc.WithDisableServiceConfig(),
	}
}

func TransportCredentialsChannelOption(secureEndpoint bool, credentialProvider auth.CredentialProvider) grpc.DialOption {
	if !secureEndpoint {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         credentialProvider.GetCacheTlsHostname(),
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(config))
}

func AllDialOptions(grpcConfig config.IGrpcConfiguration, secureEndpoint bool, credentialProvider auth.CredentialProvider, options ...grpc.DialOption) []grpc.DialOption {
	options = append(options, TransportCredentialsChannelOption(secureEndpoint, credentialProvider))
	options = append(options, GrpcChannelOptionsFromGrpcConfig(grpcConfig)...)
	return options
}
