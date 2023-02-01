package models

import (
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ControlGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type DataGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type LocalDataGrpcManagerRequest struct {
	Endpoint string
}
type CacheGetRequest struct {
	// CacheName Name of the cache to get the item from
	CacheName string
	// Key []byte key to be used to retrieve item from cache
	Key []byte
}

type CacheSetRequest struct {
	CacheName  string
	Key        []byte
	Value      []byte
	TtlSeconds uint32
}

type CacheDeleteRequest struct {
	CacheName string
	Key       []byte
}

type CreateCacheRequest struct {
	CacheName string
}

type DeleteCacheRequest struct {
	CacheName string
}

type ListCachesRequest struct {
	NextToken string
}

type TopicSubscribeRequest struct {
	CacheName string
	TopicName string
}

type TopicPublishRequest struct {
	CacheName string
	TopicName string
	Value     TopicValue
}

type TopicValue interface {
	isTopicValue()
}

type TopicValueBytes struct {
	Bytes []byte
}

func (_ TopicValueBytes) isTopicValue() {}

type TopicValueString struct {
	Text string
}

func (_ TopicValueString) isTopicValue() {}

type ControlClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

type DataClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTtlSeconds  uint32
}

type PubSubClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

type NewLocalPubSubClientRequest struct {
	Port int
}

type ConvertEcacheResultRequest struct {
	ECacheResult pb.ECacheResult
	Message      string
	OpName       string
}
