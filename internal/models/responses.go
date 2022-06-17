package models

import (
	"encoding/json"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type CreateCacheRequest struct {
	CacheName string
}

type DeleteCacheRequest struct {
	CacheName string
}

type ListCachesRequest struct {
	NextToken string
}

type ListCachesResponse struct {
	NextToken string
	Caches    []CacheInfo
}

func NewListCacheResponse(resp *pb.XListCachesResponse) *ListCachesResponse {
	var caches = []CacheInfo{}
	for _, cache := range resp.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &ListCachesResponse{NextToken: resp.NextToken, Caches: caches}
}

type CacheInfo struct {
	Name string
}

func NewCacheInfo(cache *pb.XCache) CacheInfo {
	return CacheInfo{Name: cache.CacheName}
}

type CreateSigningKeyRequest struct {
	TtlMinutes uint32
}

type CreateSigningKeyResponse struct {
	KeyId     string
	Endpoint  string
	Key       string
	ExpiresAt time.Time
}

func NewCreateSigningKeyResponse(endpoint string, resp *pb.XCreateSigningKeyResponse) (*CreateSigningKeyResponse, error) {
	var keyObj map[string]string
	err := json.Unmarshal([]byte(resp.GetKey()), &keyObj)
	if err != nil {
		return nil, err
	}
	return &CreateSigningKeyResponse{
		KeyId:     keyObj["kid"],
		Endpoint:  endpoint,
		Key:       resp.GetKey(),
		ExpiresAt: time.Unix(int64(resp.GetExpiresAt()), 0),
	}, nil
}

type RevokeSigningKeyRequest struct {
	KeyId string
}

type ListSigningKeysRequest struct {
	NextToken string
}

type ListSigningKeysResponse struct {
	NextToken   string
	SigningKeys []SigningKey
}

func NewListSigningKeysResponse(endpoint string, resp *pb.XListSigningKeysResponse) *ListSigningKeysResponse {
	var signingKeys []SigningKey
	for _, signingKey := range resp.SigningKey {
		signingKeys = append(signingKeys, NewSigningKey(endpoint, signingKey))
	}
	return &ListSigningKeysResponse{NextToken: resp.GetNextToken(), SigningKeys: signingKeys}
}

type SigningKey struct {
	KeyId     string
	Endpoint  string
	ExpiresAt time.Time
}

func NewSigningKey(endpoint string, signingKey *pb.XSigningKey) SigningKey {
	return SigningKey{
		KeyId:     signingKey.GetKeyId(),
		Endpoint:  endpoint,
		ExpiresAt: time.Unix(int64(signingKey.GetExpiresAt()), 0),
	}
}

const (
	HIT  string = "HIT"
	MISS string = "MISS"
)

type CacheGetRequest struct {
	CacheName string
	Key       interface{}
}

type GetCacheResponse struct {
	Value  []byte
	Result string
}

func NewGetCacheResponse(resp *pb.XGetResponse) (*GetCacheResponse, momentoerrors.MomentoSvcErr) {
	var result string
	if resp.Result == pb.ECacheResult_Hit {
		result = HIT
	} else if resp.Result == pb.ECacheResult_Miss {
		result = MISS
	} else {
		return nil, ConvertEcacheResult(ConvertEcacheResultRequest{
			ECacheResult: resp.Result,
			Message:      resp.Message,
			OpName:       "GET",
		})
	}
	return &GetCacheResponse{Value: resp.CacheBody, Result: result}, nil
}

type CacheSetRequest struct {
	CacheName  string
	Key        interface{}
	Value      interface{}
	TtlSeconds uint32
}

type SetCacheResponse struct {
	Value []byte
}

func NewSetCacheResponse(resp *pb.XSetResponse, value []byte) *SetCacheResponse {
	return &SetCacheResponse{Value: value}
}

type CacheDeleteRequest struct {
	CacheName string
	Key       interface{}
}

type DeleteCacheResponse struct {
	Value  []byte
	Result string
}

func NewDeleteCacheResponse(resp *pb.XGetResponse) (*DeleteCacheResponse, momentoerrors.MomentoSvcErr) {
	var result string
	if resp.Result == pb.ECacheResult_Hit {
		result = HIT
	} else if resp.Result == pb.ECacheResult_Miss {
		result = MISS
	} else {
		return nil, ConvertEcacheResult(ConvertEcacheResultRequest{
			ECacheResult: resp.Result,
			Message:      resp.Message,
			OpName:       "GET",
		})
	}
	return &DeleteCacheResponse{Value: resp.CacheBody, Result: result}, nil
}