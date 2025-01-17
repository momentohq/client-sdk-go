package grpcmanagers

import (
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"google.golang.org/grpc"
)

type AuthGrpcManager struct {
	Conn      *grpc.ClientConn
	AuthToken string
}

func NewAuthGrpcManager(request *models.AuthGrpcManagerRequest) (*AuthGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetControlEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	headerInterceptors := []grpc.UnaryClientInterceptor{
		interceptor.AddAuthHeadersInterceptor(authToken),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsControlEndpointSecure(),
			grpc.WithChainUnaryInterceptor(headerInterceptors...),
		)...,
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &AuthGrpcManager{Conn: conn, AuthToken: authToken}, nil
}

func (authManager *AuthGrpcManager) Close() momentoerrors.MomentoSvcErr {
	return momentoerrors.ConvertSvcErr(authManager.Conn.Close())
}
