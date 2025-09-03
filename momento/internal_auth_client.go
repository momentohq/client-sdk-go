package momento

import (
	"context"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	auth_responses "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
)

type authClient struct {
	grpcManager *grpcmanagers.AuthGrpcManager
	grpcClient  pb.AuthClient
}

func newAuthClient(request *models.AuthClientRequest) (*authClient, momentoerrors.MomentoSvcErr) {
	// NOTE: This is hard-coded for now but we may want to expose it via TopicConfiguration in the future,
	// as we do with some of the other clients. Defaults to keep-alive pings enabled.
	grpcConfig := config.NewStaticGrpcConfiguration(&config.GrpcConfigurationProps{})
	authManager, err := grpcmanagers.NewAuthGrpcManager(&models.AuthGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  grpcConfig,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &authClient{grpcManager: authManager, grpcClient: pb.NewAuthClient(authManager.Conn)}, nil
}

func (client *authClient) Close() {
	defer client.grpcManager.Close()
}

func (client *authClient) RefreshApiKey(ctx context.Context, request *RefreshApiKeyRequest) (auth_responses.RefreshApiKeyResponse, MomentoError) {
	grpc_request := &pb.XRefreshApiTokenRequest{
		RefreshToken: request.RefreshToken,
		ApiKey:       client.grpcManager.AuthToken,
	}

	resp, err := client.grpcClient.RefreshApiToken(ctx, grpc_request)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	jsonObject := map[string]string{
		"api_key":  resp.ApiKey,
		"endpoint": resp.Endpoint,
	}
	jsonString, err := json.Marshal(jsonObject)
	if err != nil {
		return nil, NewMomentoError(
			momentoerrors.ClientSdkError,
			"Unable to map API key to necessary form",
			err,
		)
	}
	b64Encoded := b64.StdEncoding.EncodeToString([]byte(jsonString))

	return &auth_responses.RefreshApiKeySuccess{
		ApiKey:       b64Encoded,
		RefreshToken: resp.RefreshToken,
		Endpoint:     resp.Endpoint,
		ExpiresAt:    utils.ExpiresAtFromEpoch(int64(resp.ValidUntil)),
	}, nil
}

func (client *authClient) GenerateApiKey(ctx context.Context, request *GenerateApiKeyRequest) (auth_responses.GenerateApiKeyResponse, MomentoError) {
	permissions, permsErr := permissionsFromTokenScope(request.Scope)
	if permsErr != nil {
		return nil, permsErr
	}

	grpc_request := &pb.XGenerateApiTokenRequest{
		Permissions: permissions,
		AuthToken:   client.grpcManager.AuthToken,
	}

	if request.ExpiresIn.DoesExpire() {
		grpc_request.Expiry = &pb.XGenerateApiTokenRequest_Expires_{Expires: &pb.XGenerateApiTokenRequest_Expires{ValidForSeconds: uint32(request.ExpiresIn.Seconds())}}
	} else {
		grpc_request.Expiry = &pb.XGenerateApiTokenRequest_Never_{Never: &pb.XGenerateApiTokenRequest_Never{}}
	}

	resp, err := client.grpcClient.GenerateApiToken(ctx, grpc_request)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	jsonObject := map[string]string{
		"api_key":  resp.ApiKey,
		"endpoint": resp.Endpoint,
	}
	jsonString, err := json.Marshal(jsonObject)
	if err != nil {
		return nil, NewMomentoError(
			momentoerrors.ClientSdkError,
			"Unable to map API key to necessary form",
			err,
		)
	}
	b64Encoded := b64.StdEncoding.EncodeToString([]byte(jsonString))

	return &auth_responses.GenerateApiKeySuccess{
		ApiKey:       b64Encoded,
		RefreshToken: resp.RefreshToken,
		Endpoint:     resp.Endpoint,
		ExpiresAt:    utils.ExpiresAtFromEpoch(int64(resp.ValidUntil)),
	}, nil
}

func permissionsFromTokenScope(scope PermissionScope) (*pb.Permissions, MomentoError) {
	var permissionsObject *pb.Permissions

	switch stype := scope.(type) {
	case internal.InternalSuperUserPermissions:
		permissionsObject = &pb.Permissions{
			Kind: &pb.Permissions_SuperUser{
				SuperUser: pb.SuperUserPermissions_SuperUser,
			},
		}
	case Permissions:
		var permissions []*pb.PermissionsType

		if len(stype.Permissions) == 0 {
			return nil, NewMomentoError(
				momentoerrors.InvalidArgumentError,
				"Permissions list cannot be empty",
				nil,
			)
		}

		for _, perm := range stype.Permissions {
			switch ptype := perm.(type) {
			case CachePermission:
				permissions = append(permissions, &pb.PermissionsType{
					Kind: &pb.PermissionsType_CachePermissions_{
						CachePermissions: cachePermissionsToGrpcCachePermissions(ptype),
					},
				})
				continue
			case TopicPermission:
				permissions = append(permissions, &pb.PermissionsType{
					Kind: &pb.PermissionsType_TopicPermissions_{
						TopicPermissions: topicPermissionsToGrpcTopicPermissions(ptype),
					},
				})
				continue
			}
		}

		permissionsObject = &pb.Permissions{
			Kind: &pb.Permissions_Explicit{
				Explicit: &pb.ExplicitPermissions{
					Permissions: permissions, // TODO
				},
			},
		}
	default:
		return nil, NewMomentoError(
			momentoerrors.InvalidArgumentError,
			"Unrecognized PermissionScope",
			nil,
		)
	}

	return permissionsObject, nil
}
