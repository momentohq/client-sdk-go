package momento_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMomento(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Momento Suite")
}
