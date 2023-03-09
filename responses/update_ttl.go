package responses

// UpdateTtlResponse is a base response type for an update ttl request.
type UpdateTtlResponse interface {
	isUpdateTtlResponse()
}

// UpdateTtlMiss indicates an update ttl request was not applied due to the key not being present.
type UpdateTtlMiss struct{}

func (*UpdateTtlMiss) isUpdateTtlResponse() {}

// UpdateTtlSet indicates a successful update ttl request.
type UpdateTtlSet struct{}

func (*UpdateTtlSet) isUpdateTtlResponse() {}
