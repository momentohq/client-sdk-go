package impl

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"

	"github.com/momentohq/client-sdk-go/config/compression"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
)

// GzipCompressorFactory implements the CompressionStrategyFactory interface.
type GzipCompressorFactory struct{}

func (f GzipCompressorFactory) NewCompressionStrategy(props compression.CompressionStrategyProps) compression.CompressionStrategy {
	compressionLevel := gzip.DefaultCompression
	if props.CompressionLevel == compression.CompressionLevelFastest {
		compressionLevel = gzip.BestSpeed
	} else if props.CompressionLevel == compression.CompressionLevelSmallestSize {
		compressionLevel = gzip.BestCompression
	}

	return gzipCompressor{
		compressionLevel: compressionLevel,
	}
}

// gzipCompressor implements the CompressionStrategy interface.
type gzipCompressor struct {
	compressionLevel int
}

// The byte sequence that begins a gzip compressed data frame.
// https://loc.gov/preservation/digital/formats/fdd/fdd000599.shtml#sign
const MAGIC_NUMBER = 0x1f8b

func (c gzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter, err := gzip.NewWriterLevel(&buf, c.compressionLevel)
	if err != nil {
		return nil, err
	}

	_, err = gzWriter.Write(data)
	if err != nil {
		return nil, err
	}

	err = gzWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c gzipCompressor) Decompress(data []byte, logger logger.MomentoLogger) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		logger.Error("Failed to create gzip reader", "error", err)
		return nil, err
	}
	defer reader.Close()

	if isGzipCompressed(data) {
		logger.Trace("Decompressing gzip compressed data")
		return io.ReadAll(reader)
	}
	logger.Trace("Data is not gzip compressed, passing through")
	return data, nil
}

func isGzipCompressed(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	// Extract the first 2 bytes in little endian order to compare
	// to the magic number.
	return binary.LittleEndian.Uint16(data[:2]) == MAGIC_NUMBER
}

type GzipCompressionMiddlewareProps struct {
	Logger          logger.MomentoLogger
	IncludeTypes    []interface{}
	CompressorProps compression.CompressionStrategyProps
}

// NewGzipCompressionMiddleware creates a new compression middleware that uses gzip.
// Example usage:
//
//	compressionMiddleware := gzip_compression.NewGzipCompressionMiddleware(gzip_compression.GzipCompressionMiddlewareProps{
//		Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG).GetLogger("gzip-compression-middleware"),
//		CompressionLevel: gzip_compression.CompressionLevelFastest,
//	})
func NewGzipCompressionMiddleware(props GzipCompressionMiddlewareProps) middleware.Middleware {
	compressionMiddlewareProps := compression.CompressionMiddlewareProps{
		CompressorFactory: GzipCompressorFactory{},
		CompressorProps:   props.CompressorProps,
		Logger:            props.Logger,
		IncludeTypes:      props.IncludeTypes,
	}
	return compression.NewCompressionMiddleware(compressionMiddlewareProps)
}
