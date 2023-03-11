package responses

// IncreaseTtlResponse is the base response type for an increase ttl request.
type IncreaseTtlResponse interface {
	isIncreaseTtlResponse()
}

// IncreaseTtlNotSet indicates a increase ttl request was not applied.
type IncreaseTtlNotSet struct{}

func (*IncreaseTtlNotSet) isIncreaseTtlResponse() {}

// IncreaseTtlMiss indicates a increase ttl request was not applied due to the key not being present.
type IncreaseTtlMiss struct{}

func (*IncreaseTtlMiss) isIncreaseTtlResponse() {}

// IncreaseTtlSet indicates a successful increase ttl request.
type IncreaseTtlSet struct{}

func (*IncreaseTtlSet) isIncreaseTtlResponse() {}
