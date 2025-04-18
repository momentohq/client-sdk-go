package responses

// GetWithHashResponse is the base response type for a get with hash request.
type GetWithHashResponse interface {
	MomentoCacheResponse
	isGetWithHashResponse()
}

// GetWithHashMiss Miss response to a cache GetWithHash api request.
type GetWithHashMiss struct{}

func (GetWithHashMiss) isGetWithHashResponse() {}

// GetWithHashHit Hit response to a cache GetWithHash api request.
type GetWithHashHit struct {
	value []byte
	hash  []byte
}

func (GetWithHashHit) isGetWithHashResponse() {}

// ValueString Returns value stored in cache as string for GetWithHashHit responses.
func (resp GetWithHashHit) ValueString() string {
	return string(resp.value)
}

// ValueByte Returns value stored in cache as bytes for GetWithHashHit responses.
func (resp GetWithHashHit) ValueByte() []byte {
	return resp.value
}

// HashString Returns hash of value stored in cache as string for GetWithHashHit responses.
func (resp GetWithHashHit) HashString() string {
	return string(resp.hash)
}

// HashByte Returns hash of value stored in cache as bytes for GetWithHashHit responses.
func (resp GetWithHashHit) HashByte() []byte {
	return resp.hash
}

// NewGetWithHashHit returns a new GetWithHashHit containing the supplied value.
func NewGetWithHashHit(value []byte, hash []byte) *GetWithHashHit {
	return &GetWithHashHit{value: value, hash: hash}
}
