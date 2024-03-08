package batchutils

import (
	"time"

	"github.com/momentohq/client-sdk-go/momento"
)

const defaultRequestTimeout = 10 * time.Second

func getRequestTimeout(propsTimeout *time.Duration) (requestTimeout time.Duration) {
	if propsTimeout == nil {
		requestTimeout = defaultRequestTimeout
	} else {
		requestTimeout = *propsTimeout
	}
	return
}

type errKeyVal struct {
	key   momento.Value
	error error
}
