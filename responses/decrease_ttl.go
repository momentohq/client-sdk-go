package responses

type DecreaseTtlResponse interface {
	isDecreaseTtlResponse()
}

type DecreaseTtlNotSet struct{}

func (*DecreaseTtlNotSet) isDecreaseTtlResponse() {}

type DecreaseTtlMiss struct{}

func (*DecreaseTtlMiss) isDecreaseTtlResponse() {}

type DecreaseTtlSet struct{}

func (*DecreaseTtlSet) isDecreaseTtlResponse() {}
