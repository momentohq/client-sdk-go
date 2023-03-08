package responses

type UpdateTtlResponse interface {
	isUpdateTtlResponse()
}

type UpdateTtlNotSet struct{}

func (*UpdateTtlNotSet) isUpdateTtlResponse() {}

type UpdateTtlMiss struct{}

func (*UpdateTtlMiss) isUpdateTtlResponse() {}

type UpdateTtlSet struct{}

func (*UpdateTtlSet) isUpdateTtlResponse() {}
