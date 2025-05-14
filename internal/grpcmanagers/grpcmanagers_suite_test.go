package grpcmanagers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGrpcManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GrpcManagers Suite")
}
