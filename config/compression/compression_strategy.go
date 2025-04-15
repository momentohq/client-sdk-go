package compression

import "github.com/momentohq/client-sdk-go/config/logger"

type CompressionStrategy interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

type CompressionStrategyProps struct {
	CompressionLevel CompressionLevel
	Logger           logger.MomentoLogger
}

type CompressionLevel string

const (
	CompressionLevelDefault      CompressionLevel = "default"
	CompressionLevelFastest      CompressionLevel = "fastest"
	CompressionLevelSmallestSize CompressionLevel = "smallestSize"
)

type CompressionStrategyFactory interface {
	NewCompressionStrategy(props CompressionStrategyProps) CompressionStrategy
}
