package responses

// GetResponse is the base response type for a get request.
type GetResponse interface {
	isGetResponse()
}

// GetMiss Miss response to a cache Get api request.
type GetMiss struct{}

func (GetMiss) isGetResponse() {}

// GetHit Hit response to a cache Get api request.
type GetHit struct {
	value []byte
}

func (GetHit) isGetResponse() {}

// ValueString Returns value stored in cache as string if there was Hit. Returns an empty string otherwise.
func (resp GetHit) ValueString() string {
	return string(resp.value)
}

// ValueByte Returns value stored in cache as bytes if there was Hit. Returns nil otherwise.
func (resp GetHit) ValueByte() []byte {
	return resp.value
}

// NewGetHit returns a new GetHit containing the supplied value.
func NewGetHit(value []byte) *GetHit {
	return &GetHit{value: value}
}
