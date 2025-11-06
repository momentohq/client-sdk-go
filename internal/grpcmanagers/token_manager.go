package grpcmanagers

import (
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
)

type TokenGrpcManager struct {
	Conn      *grpc.ClientConn
	AuthToken string
}

func NewTokenGrpcManager(request *models.TokenGrpcManagerRequest) (*TokenGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetTokenEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	headerInterceptors := []grpc.UnaryClientInterceptor{
		interceptor.AddAuthHeadersInterceptor(authToken),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsTokenEndpointSecure(),
			request.CredentialProvider,
			grpc.WithChainUnaryInterceptor(headerInterceptors...),
		)...,
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &TokenGrpcManager{Conn: conn, AuthToken: authToken}, nil
}

func (tokenManager *TokenGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(tokenManager.Conn.Close())
}
