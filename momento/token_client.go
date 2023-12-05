package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	auth_responses "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
)

type tokenClient struct {
	grpcManager *grpcmanagers.TokenGrpcManager
	grpcClient  pb.TokenClient
}

func newTokenClient(request *models.TokenClientRequest) (*tokenClient, momentoerrors.MomentoSvcErr) {
	tokenManager, err := grpcmanagers.NewTokenGrpcManager(&models.TokenGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &tokenClient{grpcManager: tokenManager, grpcClient: pb.NewTokenClient(tokenManager.Conn)}, nil
}

func (client *tokenClient) close() {
	defer client.grpcManager.Close()
}

func cacheRoleToGrpcCacheRole(role CacheRole) pb.CacheRole {
	switch role {
	case ReadWrite:
		return pb.CacheRole_CacheReadWrite
	case ReadOnly:
		return pb.CacheRole_CacheReadOnly
	case WriteOnly:
		return pb.CacheRole_CacheWriteOnly
	}
	return pb.CacheRole_CachePermitNone
}

func topicRoleToGrpcTopicRole(role TopicRole) pb.TopicRole {
	switch role {
	case PublishOnly:
		return pb.TopicRole_TopicWriteOnly
	case SubscribeOnly:
		return pb.TopicRole_TopicReadOnly
	case PublishSubscribe:
		return pb.TopicRole_TopicReadWrite
	}
	return pb.TopicRole_TopicPermitNone
}

func cachePermissionsToGrpcCachePermissions(cp CachePermission) *pb.PermissionsType_CachePermissions {
	switch cp.Cache.(type) {
	case AllCaches:
		return &pb.PermissionsType_CachePermissions{
			Cache: &pb.PermissionsType_CachePermissions_AllCaches{},
			Role:  cacheRoleToGrpcCacheRole(cp.Role),
		}
	case CacheName:
		return &pb.PermissionsType_CachePermissions{
			Cache: &pb.PermissionsType_CachePermissions_CacheSelector{
				CacheSelector: &pb.PermissionsType_CacheSelector{
					Kind: &pb.PermissionsType_CacheSelector_CacheName{CacheName: cp.Cache.CacheName()},
				},
			},
			Role: cacheRoleToGrpcCacheRole(cp.Role),
		}
	}

	return nil
}

func topicPermissionsToGrpcTopicPermissions(tp TopicPermission) *pb.PermissionsType_TopicPermissions {
	topicPermissions := &pb.PermissionsType_TopicPermissions{
		Role: topicRoleToGrpcTopicRole(tp.Role),
	}
	switch tp.Cache.(type) {
	case AllCaches:
		topicPermissions.Cache = &pb.PermissionsType_TopicPermissions_AllCaches{}
		switch tp.Topic.(type) {
		case AllTopics:
			topicPermissions.Topic = &pb.PermissionsType_TopicPermissions_AllTopics{}
			return topicPermissions
		case TopicNamePrefix:
			topicPermissions.Topic = &pb.PermissionsType_TopicPermissions_TopicSelector{
				TopicSelector: &pb.PermissionsType_TopicSelector{
					Kind: &pb.PermissionsType_TopicSelector_TopicNamePrefix{
						TopicNamePrefix: tp.Topic.TopicName(),
					},
				},
			}
			return topicPermissions
		default:
			topicPermissions.Topic = &pb.PermissionsType_TopicPermissions_TopicSelector{
				TopicSelector: &pb.PermissionsType_TopicSelector{
					Kind: &pb.PermissionsType_TopicSelector_TopicName{
						TopicName: tp.Topic.TopicName(),
					},
				},
			}
			return topicPermissions
		}

	case CacheName:
		topicPermissions.Cache = &pb.PermissionsType_TopicPermissions_CacheSelector{
			CacheSelector: &pb.PermissionsType_CacheSelector{
				Kind: &pb.PermissionsType_CacheSelector_CacheName{CacheName: tp.Cache.CacheName()},
			},
		}
		switch tp.Topic.(type) {
		case AllTopics:
			topicPermissions.Topic = &pb.PermissionsType_TopicPermissions_AllTopics{}
			return topicPermissions
		case TopicNamePrefix:
			topicPermissions.Topic = &pb.PermissionsType_TopicPermissions_TopicSelector{
				TopicSelector: &pb.PermissionsType_TopicSelector{
					Kind: &pb.PermissionsType_TopicSelector_TopicNamePrefix{
						TopicNamePrefix: tp.Topic.TopicName(),
					},
				},
			}
			return topicPermissions
		default:
			topicPermissions.Topic = &pb.PermissionsType_TopicPermissions_TopicSelector{
				TopicSelector: &pb.PermissionsType_TopicSelector{
					Kind: &pb.PermissionsType_TopicSelector_TopicName{
						TopicName: tp.Topic.TopicName(),
					},
				},
			}
			return topicPermissions
		}
	}

	return nil
}

func disposableTokenPermissionsToGrpcDisposablePermissions(cp DisposableTokenCachePermission) *pb.PermissionsType_CachePermissions {
	switch cp.Cache.(type) {
	case AllCaches:
		switch itype := cp.Item.(type) {
		case AllCacheItems:
			return &pb.PermissionsType_CachePermissions{
				Cache:     &pb.PermissionsType_CachePermissions_AllCaches{},
				Role:      cacheRoleToGrpcCacheRole(cp.Role),
				CacheItem: &pb.PermissionsType_CachePermissions_AllItems{},
			}
		case CacheItemKey:
			return &pb.PermissionsType_CachePermissions{
				Cache: &pb.PermissionsType_CachePermissions_AllCaches{},
				Role:  cacheRoleToGrpcCacheRole(cp.Role),
				CacheItem: &pb.PermissionsType_CachePermissions_ItemSelector{
					ItemSelector: &pb.PermissionsType_CacheItemSelector{
						Kind: &pb.PermissionsType_CacheItemSelector_Key{
							Key: itype.Key,
						},
					},
				},
			}
		case CacheItemKeyPrefix:
			return &pb.PermissionsType_CachePermissions{
				Cache: &pb.PermissionsType_CachePermissions_AllCaches{},
				Role:  cacheRoleToGrpcCacheRole(cp.Role),
				CacheItem: &pb.PermissionsType_CachePermissions_ItemSelector{
					ItemSelector: &pb.PermissionsType_CacheItemSelector{
						Kind: &pb.PermissionsType_CacheItemSelector_KeyPrefix{
							KeyPrefix: itype.KeyPrefix,
						},
					},
				},
			}
		}
	case CacheName:
		cacheSelector := &pb.PermissionsType_CachePermissions_CacheSelector{
			CacheSelector: &pb.PermissionsType_CacheSelector{
				Kind: &pb.PermissionsType_CacheSelector_CacheName{CacheName: cp.Cache.CacheName()},
			},
		}
		switch itype := cp.Item.(type) {
		case AllCacheItems:
			return &pb.PermissionsType_CachePermissions{
				Cache:     cacheSelector,
				Role:      cacheRoleToGrpcCacheRole(cp.Role),
				CacheItem: &pb.PermissionsType_CachePermissions_AllItems{},
			}
		case CacheItemKey:
			return &pb.PermissionsType_CachePermissions{
				Cache: cacheSelector,
				Role:  cacheRoleToGrpcCacheRole(cp.Role),
				CacheItem: &pb.PermissionsType_CachePermissions_ItemSelector{
					ItemSelector: &pb.PermissionsType_CacheItemSelector{
						Kind: &pb.PermissionsType_CacheItemSelector_Key{
							Key: itype.Key,
						},
					},
				},
			}
		case CacheItemKeyPrefix:
			return &pb.PermissionsType_CachePermissions{
				Cache: cacheSelector,
				Role:  cacheRoleToGrpcCacheRole(cp.Role),
				CacheItem: &pb.PermissionsType_CachePermissions_ItemSelector{
					ItemSelector: &pb.PermissionsType_CacheItemSelector{
						Kind: &pb.PermissionsType_CacheItemSelector_KeyPrefix{
							KeyPrefix: itype.KeyPrefix,
						},
					},
				},
			}
		}
	}

	return nil
}

func (client *tokenClient) GenerateDisposableToken(ctx context.Context, request *GenerateDisposableTokenRequest) (auth_responses.GenerateDisposableTokenResponse, error) {
	var permissions []*pb.PermissionsType
	switch stype := request.Scope.(type) {
	case Permissions:
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
				// do topic permission things
			case DisposableTokenCachePermission:
				permissions = append(permissions, &pb.PermissionsType{
					Kind: &pb.PermissionsType_CachePermissions_{
						CachePermissions: disposableTokenPermissionsToGrpcDisposablePermissions(ptype),
					},
				})
				continue
			}
		}
	}
	if err := utils.ValidateDisposableTokenExpiry(request.ExpiresIn); err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	resp, err := client.grpcClient.GenerateDisposableToken(ctx, &pb.XGenerateDisposableTokenRequest{
		AuthToken: client.grpcManager.AuthToken,
		Expires: &pb.XGenerateDisposableTokenRequest_Expires{
			ValidForSeconds: uint32(request.ExpiresIn.Seconds()),
		},
		Permissions: &pb.Permissions{
			Kind: &pb.Permissions_Explicit{
				Explicit: &pb.ExplicitPermissions{
					Permissions: permissions,
				},
			},
		},
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &auth_responses.GenerateDisposableTokenSuccess{
		ApiKey:     resp.ApiKey,
		Endpoint:   resp.Endpoint,
		ValidUntil: resp.ValidUntil,
	}, nil
}
