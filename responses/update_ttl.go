package responses

type UpdateTtlResponse interface {
	isUpdateTtlResponse()
}

type UpdateTtlMiss struct{}

func (*UpdateTtlMiss) isUpdateTtlResponse() {}

type UpdateTtlSet struct{}

func (*UpdateTtlSet) isUpdateTtlResponse() {}
