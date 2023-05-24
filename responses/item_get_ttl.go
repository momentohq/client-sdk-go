package responses

import "time"

// ItemGetTtlResponse is the base type for item get TTL requests.
type ItemGetTtlResponse interface {
	isItemGetTtlResponse()
}

// ItemGetTtlHit hit response to an item get TTL request.
type ItemGetTtlHit struct {
	value uint64
}

func (r *ItemGetTtlHit) isItemGetTtlResponse() {}

// Ttl returns the TTL in milliseconds
func (r *ItemGetTtlHit) Ttl() uint64 {
	return r.value
}

// TtlDuration returns the TTL as a duration
func (r *ItemGetTtlHit) TtlDuration() time.Duration {
	return time.Millisecond * time.Duration(r.value)
}

// NewItemGetTtlHit returns a new ItemGetTtlHit containing the item TTL.
func NewItemGetTtlHit(value uint64) *ItemGetTtlHit {
	return &ItemGetTtlHit{value: value}
}

type ItemGetTtlMiss struct{}

func (r *ItemGetTtlMiss) isItemGetTtlResponse() {}
