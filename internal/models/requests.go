package models

import (
	"time"

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

func (TopicValueBytes) isTopicValue() {}

type TopicValueString struct {
	Text string
}

func (TopicValueString) isTopicValue() {}

type ControlClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

type DataClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTtl         time.Duration
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
