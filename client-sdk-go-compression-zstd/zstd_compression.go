package zstd_compression

import (
	"fmt"

	"github.com/klauspost/compress/zstd"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

type zstdCompressionMiddleware struct {
	middleware.Middleware
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

type ZstdCompressionMiddlewareProps struct {
	Logger           logger.MomentoLogger
	IncludeTypes     []interface{}
	CompressionLevel zstd.EncoderLevel
}

func (mw *zstdCompressionMiddleware) GetRequestHandler(baseHandler middleware.RequestHandler) (middleware.RequestHandler, error) {
	return NewZstdCompressionMiddlewareRequestHandler(baseHandler, mw.encoder, mw.decoder), nil
}

func NewZstdCompressionMiddleware(props ZstdCompressionMiddlewareProps) middleware.Middleware {
	compressionLevel := zstd.SpeedDefault
	if props.CompressionLevel != 0 {
		compressionLevel = props.CompressionLevel
	}
	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(compressionLevel))
	decoder, _ := zstd.NewReader(nil)
	mw := middleware.NewMiddleware(middleware.Props{
		Logger:       props.Logger,
		IncludeTypes: props.IncludeTypes,
	})
	return &zstdCompressionMiddleware{mw, encoder, decoder}
}

type zstdCompressionMiddlewareRequestHandler struct {
	middleware.RequestHandler
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewZstdCompressionMiddlewareRequestHandler(rh middleware.RequestHandler, encoder *zstd.Encoder, decoder *zstd.Decoder) middleware.RequestHandler {
	return &zstdCompressionMiddlewareRequestHandler{rh, encoder, decoder}
}

func (rh *zstdCompressionMiddlewareRequestHandler) compress(requestType string, rawData []byte) []byte {
	compressed := rh.encoder.EncodeAll(rawData, nil)
	rh.GetLogger().Info("Compressed request %s: %d bytes -> %d bytes", requestType, len(rawData), len(compressed))
	return compressed
}

func (rh *zstdCompressionMiddlewareRequestHandler) decompress(responseType string, rawData []byte) ([]byte, error) {
	decompressed, err := rh.decoder.DecodeAll(rawData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress response: %v", err)
	}
	rh.GetLogger().Info("Decompressed response %s: %d bytes -> %d bytes", responseType, len(rawData), len(decompressed))
	return decompressed, nil
}

// We currently compress only on these scalar write requests:
// Set, SetIfAbsent, SetIfPresent, SetIfEqual, SetIfNotEqual, SetIfAbsentOrEqual, SetIfPresentAndNotEqual,
// SetWithHash, SetIfPresentAndHashEqual, SetIfPresentAndHashNotEqual, SetIfAbsentOrHashEqual, SetIfAbsentOrHashNotEqual.
// Specify IncludeTypes in ZstdCompressionMiddlewareProps if you wish to compress only a subset of these requests.
func (rh *zstdCompressionMiddlewareRequestHandler) OnRequest(req interface{}) (interface{}, error) {
	// We still need to use a switch statement to be able to access the request objects
	// as the specific request types in order to access the Value field.
	switch r := req.(type) {
	case *momento.SetRequest:
		compressed := rh.compress("SetRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentRequest:
		compressed := rh.compress("SetIfAbsentRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentRequest:
		compressed := rh.compress("SetIfPresentRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfEqualRequest:
		compressed := rh.compress("SetIfEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfNotEqualRequest:
		compressed := rh.compress("SetIfNotEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrEqualRequest:
		compressed := rh.compress("SetIfAbsentOrEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndNotEqualRequest:
		compressed := rh.compress("SetIfPresentAndNotEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetWithHashRequest:
		compressed := rh.compress("SetWithHashRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndHashEqualRequest:
		compressed := rh.compress("SetIfPresentAndHashEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndHashNotEqualRequest:
		compressed := rh.compress("SetIfPresentAndHashNotEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrHashEqualRequest:
		compressed := rh.compress("SetIfAbsentOrHashEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrHashNotEqualRequest:
		compressed := rh.compress("SetIfAbsentOrHashNotEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	default:
		rh.GetLogger().Info("No action for OnRequest type: %T", req)
		return req, nil
	}
}

// We currently decompress only on these scalar read responses: Get, GetWithHash.
// Specify IncludeTypes in ZstdCompressionMiddlewareProps if you wish to decompress only a subset of these responses.
func (rh *zstdCompressionMiddlewareRequestHandler) OnResponse(resp interface{}) (interface{}, error) {
	// We still need to use a switch statement to be able to access the response objects
	// as the specific response types in order to access the Value field.
	switch r := resp.(type) {
	case *responses.GetHit:
		decompressed, err := rh.decompress("GetHit", r.ValueByte())
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		return responses.NewGetHit(decompressed), nil
	case *responses.GetWithHashHit:
		decompressed, err := rh.decompress("GetWithHashHit", r.ValueByte())
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		return responses.NewGetWithHashHit(decompressed, r.HashByte()), nil
	default:
		rh.GetLogger().Info("No action for OnResponse type: %T", resp)
		return resp, nil
	}
}
