package zstd_compression

import (
	"github.com/klauspost/compress/zstd"
	"github.com/momentohq/client-sdk-go/config/compression"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
)

// ZstdCompressorFactory implements the CompressionStrategyFactory interface.

type ZstdCompressorFactory struct{}

func (f ZstdCompressorFactory) NewCompressionStrategy(props compression.CompressionStrategyProps) compression.CompressionStrategy {
	compressionLevel := zstd.SpeedDefault
	if props.CompressionLevel == compression.CompressionLevelFastest {
		compressionLevel = zstd.SpeedFastest
	} else if props.CompressionLevel == compression.CompressionLevelSmallestSize {
		compressionLevel = zstd.SpeedBestCompression
	}

	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(compressionLevel))
	decoder, _ := zstd.NewReader(nil)
	return zstdCompressor{
		encoder: encoder,
		decoder: decoder,
	}
}

// zstdCompressor implements the CompressionStrategy interface.

type zstdCompressor struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func (c zstdCompressor) Compress(data []byte) ([]byte, error) {
	return c.encoder.EncodeAll(data, nil), nil
}

func (c zstdCompressor) Decompress(data []byte) ([]byte, error) {
	return c.decoder.DecodeAll(data, nil)
}

type ZstdCompressionMiddlewareProps struct {
	Logger          logger.MomentoLogger
	IncludeTypes    []interface{}
	CompressorProps compression.CompressionStrategyProps
}

// NewZstdCompressionMiddleware creates a new compression middleware that uses zstd.
// Example usage:
//
//	compressionMiddleware := zstd_compression.NewZstdCompressionMiddleware(zstd_compression.ZstdCompressionMiddlewareProps{
//		Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG).GetLogger("zstd-compression-middleware"),
//		CompressionLevel: zstd_compression.CompressionLevelFastest,
//	})
func NewZstdCompressionMiddleware(props ZstdCompressionMiddlewareProps) middleware.Middleware {
	compressionMiddlewareProps := compression.CompressionMiddlewareProps{
		CompressorFactory: ZstdCompressorFactory{},
		CompressorProps:   props.CompressorProps,
		Logger:            props.Logger,
		IncludeTypes:      props.IncludeTypes,
	}
	return compression.NewCompressionMiddleware(compressionMiddlewareProps)
}
