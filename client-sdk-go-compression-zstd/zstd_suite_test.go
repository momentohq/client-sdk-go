package zstd_compression_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestZstdCompressionMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zstd Compression Middleware Suite")
}
