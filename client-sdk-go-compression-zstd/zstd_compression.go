package zstd_compression

import (
	"encoding/binary"
	"fmt"

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

	if props.Logger == nil {
		props.Logger = logger.NewNoopMomentoLoggerFactory().GetLogger("zstd-compression")
	}

	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(compressionLevel))
	decoder, _ := zstd.NewReader(nil)
	return zstdCompressor{
		encoder: encoder,
		decoder: decoder,
		logger:  props.Logger,
	}
}

// zstdCompressor implements the CompressionStrategy interface.

type zstdCompressor struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
	logger  logger.MomentoLogger
}

// The byte sequence that begins a ZSTD compressed data frame.
// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md
const MAGIC_NUMBER = 0xfd2fb528

func (c zstdCompressor) Compress(data []byte) ([]byte, error) {
	compressed := c.encoder.EncodeAll(data, nil)
	c.logger.Trace("Compressed request: %d bytes -> %d bytes", len(data), len(compressed))
	return compressed, nil
}

func (c zstdCompressor) Decompress(data []byte) ([]byte, error) {
	if isZstdCompressed(data) {
		decompressed, err := c.decoder.DecodeAll(data, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		c.logger.Trace("Decompressed response: %d bytes -> %d bytes", len(data), len(decompressed))
		return decompressed, nil
	}
	c.logger.Trace("Data is not ZSTD compressed, passing through")
	return data, nil
}

func isZstdCompressed(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	// Extract the first 4 bytes in little endian order to compare
	// to the magic number.
	return binary.LittleEndian.Uint32(data[:4]) == MAGIC_NUMBER
}

type ZstdCompressionMiddlewareProps struct {
	IncludeTypes             []interface{}
	CompressionStrategyProps compression.CompressionStrategyProps
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
		CompressorFactory:        ZstdCompressorFactory{},
		CompressionStrategyProps: props.CompressionStrategyProps,
		IncludeTypes:             props.IncludeTypes,
	}
	return compression.NewCompressionMiddleware(compressionMiddlewareProps)
}
