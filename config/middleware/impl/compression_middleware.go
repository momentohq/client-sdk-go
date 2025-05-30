package impl

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config/compression"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

// CompressionMiddleware is the base type for a middleware that compresses and decompresses
// scalar get and set requests and responses.
type CompressionMiddleware struct {
	middleware.Middleware
	compressor compression.CompressionStrategy
}

type CompressionMiddlewareProps struct {
	IncludeTypes             []interface{}
	CompressionStrategyProps compression.CompressionStrategyProps
	CompressorFactory        compression.CompressionStrategyFactory
}

func NewCompressionMiddleware(props CompressionMiddlewareProps) middleware.Middleware {
	mw := middleware.NewMiddleware(middleware.Props{
		Logger:       props.CompressionStrategyProps.Logger,
		IncludeTypes: props.IncludeTypes,
	})
	compressor := props.CompressorFactory.NewCompressionStrategy(props.CompressionStrategyProps)
	return &CompressionMiddleware{
		Middleware: mw,
		compressor: compressor,
	}
}

type CompressionMiddlewareRequestHandler struct {
	middleware.RequestHandler
	compressor compression.CompressionStrategy
}

func NewCompressionMiddlewareRequestHandler(rh middleware.RequestHandler, compressor compression.CompressionStrategy) middleware.RequestHandler {
	return &CompressionMiddlewareRequestHandler{rh, compressor}
}

func (mw *CompressionMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewCompressionMiddlewareRequestHandler(
		baseHandler,
		mw.compressor,
	), nil
}

// getValueBytes converts a momento.Value to []byte
func getValueBytes(value momento.Value) ([]byte, error) {
	switch v := value.(type) {
	case momento.Bytes:
		return v, nil
	case momento.String:
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("unsupported value type: %T", v)
	}
}

// We currently compress only on these scalar write requests:
// Set, SetIfAbsent, SetIfPresent, SetIfEqual, SetIfNotEqual, SetIfAbsentOrEqual, SetIfPresentAndNotEqual,
// SetWithHash, SetIfPresentAndHashEqual, SetIfPresentAndHashNotEqual, SetIfAbsentOrHashEqual, SetIfAbsentOrHashNotEqual.
// Specify IncludeTypes in CompressionMiddlewareProps if you wish to compress only a subset of these requests.
func (rh *CompressionMiddlewareRequestHandler) OnRequest(req interface{}) (interface{}, error) {
	// We still need to use a switch statement to be able to access the request objects
	// as the specific request types in order to access the Value field.
	switch r := req.(type) {
	case *momento.SetRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfAbsentRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfPresentRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfNotEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfNotEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfAbsentOrEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndNotEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfPresentAndNotEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetWithHashRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetWithHashRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndHashEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfPresentAndHashEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndHashNotEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfPresentAndHashNotEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrHashEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfAbsentOrHashEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrHashNotEqualRequest:
		rawData, err := getValueBytes(r.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get value bytes: %v", err)
		}
		compressed, err := rh.compressor.Compress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress SetIfAbsentOrHashNotEqualRequest: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	default:
		rh.GetLogger().Info("No action for OnRequest type: %T", req)
		return req, nil
	}
}

// We currently decompress only on these scalar read responses: Get, GetWithHash.
// Specify IncludeTypes in CompressionMiddlewareProps if you wish to decompress only a subset of these responses.
func (rh *CompressionMiddlewareRequestHandler) OnResponse(resp interface{}) (interface{}, error) {
	// We still need to use a switch statement to be able to access the response objects
	// as the specific response types in order to access the Value field.
	switch r := resp.(type) {
	case *responses.GetHit:
		rawData := r.ValueByte()
		decompressed, err := rh.compressor.Decompress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress GetHit response: %v", err)
		}
		return responses.NewGetHit(decompressed), nil
	case *responses.GetWithHashHit:
		rawData := r.ValueByte()
		decompressed, err := rh.compressor.Decompress(rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress GetWithHashHit response: %v", err)
		}
		return responses.NewGetWithHashHit(decompressed, r.HashByte()), nil
	default:
		rh.GetLogger().Info("No action for OnResponse type: %T", resp)
		return resp, nil
	}
}
