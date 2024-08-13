package helpers

import (
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/momento"
)

func NewStringKey() momento.String {
	return momento.String(uuid.NewString())
}

func NewByteKey() momento.Bytes {
	return momento.Bytes([]byte(uuid.NewString()))
}
