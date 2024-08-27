package helpers

import (
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/momento"
)

func NewRandomString() string {
	return uuid.NewString()
}

func NewRandomMomentoString() momento.String {
	return momento.String(NewRandomString())
}

func NewRandomMomentoBytes() momento.Bytes {
	return momento.Bytes([]byte(uuid.NewString()))
}
