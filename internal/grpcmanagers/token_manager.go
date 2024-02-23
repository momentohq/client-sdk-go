package grpcmanagers

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"google.golang.org/grpc"
)

type TokenGrpcManager struct {
	Conn      *grpc.ClientConn
	AuthToken string
}

const TokenPort = ":443"

func NewTokenGrpcManager(request *models.TokenGrpcManagerRequest) (*TokenGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := fmt.Sprint(request.CredentialProvider.GetTokenEndpoint(), TokenPort)
	authToken := request.CredentialProvider.GetAuthToken()
	conn, err := grpc.Dial(
		endpoint,
		AllDialOptions(
			nil,
			grpc.WithUnaryInterceptor(interceptor.AddAuthHeadersInterceptor(authToken)),
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
