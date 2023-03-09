package responses

type IncreaseTtlResponse interface {
	isIncreaseTtlResponse()
}

type IncreaseTtlNotSet struct{}

func (*IncreaseTtlNotSet) isIncreaseTtlResponse() {}

type IncreaseTtlMiss struct{}

func (*IncreaseTtlMiss) isIncreaseTtlResponse() {}

type IncreaseTtlSet struct{}

func (*IncreaseTtlSet) isIncreaseTtlResponse() {}
