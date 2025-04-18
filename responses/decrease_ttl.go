package responses

// DecreaseTtlResponse is the base response type for a decrease ttl request.
type DecreaseTtlResponse interface {
	MomentoCacheResponse
	isDecreaseTtlResponse()
}

// DecreaseTtlNotSet indicates a decrease ttl request was not applied.
type DecreaseTtlNotSet struct{}

func (*DecreaseTtlNotSet) isDecreaseTtlResponse() {}

// DecreaseTtlMiss indicates a decrease ttl request was not applied due to the key not being present.
type DecreaseTtlMiss struct{}

func (*DecreaseTtlMiss) isDecreaseTtlResponse() {}

// DecreaseTtlSet indicates a successful decrease ttl request.
type DecreaseTtlSet struct{}

func (*DecreaseTtlSet) isDecreaseTtlResponse() {}
