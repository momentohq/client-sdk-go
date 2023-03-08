package responses

type IncreaseTtlResponse interface {
	isUpdateTtlResponse()
}

type IncreaseTtlNotSet struct{}

func (*IncreaseTtlNotSet) isUpdateTtlResponse() {}

type IncreaseTtlMiss struct{}

func (*IncreaseTtlMiss) isUpdateTtlResponse() {}

type IncreaseTtlSet struct{}

func (*IncreaseTtlSet) isUpdateTtlResponse() {}
