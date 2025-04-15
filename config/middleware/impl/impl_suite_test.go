package impl_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMiddlewareImpl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Impl Suite")
}
