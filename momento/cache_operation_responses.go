package momento

import "time"

// Output of the List caches operation.
type ListCachesResponse struct {
	nextToken string
	caches    []CacheInfo
}

// Next Page Token returned by Simple Cache Service along with the list of caches.
// If nextToken is present, then this token must be provided in the next call to continue paginating through the list.
// This is done by setting this value in ListCachesRequest.
func (resp *ListCachesResponse) NextToken() string {
	return resp.nextToken
}

// Returns all caches.
func (resp *ListCachesResponse) Caches() []CacheInfo {
	return resp.caches
}

// Output of the Create Signing Key operationn
type CreateSigningKeyResponse struct {
	keyId     string
	endpoint  string
	key       string
	expiresAt time.Time
}

// Returns the Momento signing key's ID
func (resp *CreateSigningKeyResponse) KeyId() string {
	return resp.keyId
}

// Returns the Momento signing key's endpoint
func (resp *CreateSigningKeyResponse) Endpoint() string {
	return resp.endpoint
}

// Returns the Momento signing key's metadata as a JSON string
func (resp *CreateSigningKeyResponse) Key() string {
	return resp.key
}

// Returns the Momento signing key's time in which it expires
func (resp *CreateSigningKeyResponse) ExpiresAt() time.Time {
	return resp.expiresAt
}

// Output of the List Signing Keys operation
type ListSigningKeysResponse struct {
	nextToken   string
	signingKeys []SigningKey
}

// Next Page Token returned by Simple Cache Service along with the list of Momento signing keys.
// If nextToken is present, then this token must be provided in the next call to continue paginating through the list.
// This is done by setting this value in ListSigningKeysRequest.
func (resp *ListSigningKeysResponse) NextToken() string {
	return resp.nextToken
}

// Returns all Momento signing keys
func (resp *ListSigningKeysResponse) SigningKeys() []SigningKey {
	return resp.signingKeys
}

// Information about the Signing Key
type SigningKey struct {
	keyId     string
	endpoint  string
	expiresAt time.Time
}

// Returns the Momento signing key's ID
func (sk SigningKey) KeyId() string {
	return sk.keyId
}

// Returns the Momento signing key's endpoint
func (sk SigningKey) Endpoint() string {
	return sk.endpoint
}

// Returns the Momento signing key's time in which it expires
func (sk SigningKey) ExpiresAt() time.Time {
	return sk.expiresAt
}

// Information about the Cache.
type CacheInfo struct {
	name string
}

// Returns cache's name.
func (ci CacheInfo) Name() string {
	return ci.name
}

const (
	// Represents cache hit.
	HIT string = "HIT"
	// Represents cache miss.
	MISS string = "MISS"
)

// Initializes GetCacheResponse to handle gRPC get response.
type GetCacheResponse struct {
	value  []byte
	result string
}

// Returns value stored in cache as string if there was Hit. Returns an empty string otherwise.
func (resp *GetCacheResponse) StringValue() string {
	if resp.result == HIT {
		return string(resp.value)
	}
	return ""
}

// Returns value stored in cache as bytes if there was Hit. Returns nil otherwise.
func (resp *GetCacheResponse) ByteValue() []byte {
	if resp.result == HIT {
		return resp.value
	}
	return nil
}

// Returns get operation result such as HIT or MISS.
func (resp *GetCacheResponse) Result() string {
	return resp.result
}

// Initializes SetCacheResponse to handle gRPC set response.
type SetCacheResponse struct {
	value []byte
}

// Decodes and returns byte value set in cache to string.
func (resp *SetCacheResponse) StringValue() string {
	return string(resp.value)
}

// Returns byte value set in cache.
func (resp *SetCacheResponse) ByteValue() []byte {
	return resp.value
}

// Initializes DeleteCacheResponse to handle gRPC get response.
type DeleteCacheResponse struct {
	value  []byte
	result string
}

// Returns value stored in cache as string if there was Hit. Returns an empty string otherwise.
func (resp *DeleteCacheResponse) StringValue() string {
	if resp.result == HIT {
		return string(resp.value)
	}
	return ""
}

// Returns value stored in cache as bytes if there was Hit. Returns nil otherwise.
func (resp *DeleteCacheResponse) ByteValue() []byte {
	if resp.result == HIT {
		return resp.value
	}
	return nil
}

// Returns get operation result such as HIT or MISS.
func (resp *DeleteCacheResponse) Result() string {
	return resp.result
}