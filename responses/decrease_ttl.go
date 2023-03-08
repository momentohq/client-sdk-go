package responses

type DecreaseTtlResponse interface {
	isUpdateTtlResponse()
}

type DecreaseTtlNotSet struct{}

func (*DecreaseTtlNotSet) isUpdateTtlResponse() {}

type DecreaseTtlMiss struct{}

func (*DecreaseTtlMiss) isUpdateTtlResponse() {}

type DecreaseTtlSet struct{}

func (*DecreaseTtlSet) isUpdateTtlResponse() {}
