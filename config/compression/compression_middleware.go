package compression

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

type CompressionMiddleware struct {
	middleware.Middleware
	compressor CompressionStrategy
}

type CompressionMiddlewareProps struct {
	Logger            logger.MomentoLogger
	IncludeTypes      []interface{}
	CompressorProps   CompressionStrategyProps
	CompressorFactory CompressionStrategyFactory
}

func NewCompressionMiddleware(props CompressionMiddlewareProps) middleware.Middleware {
	mw := middleware.NewMiddleware(middleware.Props{
		Logger:       props.Logger,
		IncludeTypes: props.IncludeTypes,
	})
	compressor := props.CompressorFactory.NewCompressionStrategy(props.CompressorProps)
	return &CompressionMiddleware{
		Middleware: mw,
		compressor: compressor,
	}
}

type CompressionMiddlewareRequestHandler struct {
	middleware.RequestHandler
	compressor CompressionStrategy
}

func NewCompressionMiddlewareRequestHandler(rh middleware.RequestHandler, compressor CompressionStrategy) middleware.RequestHandler {
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

// We currently compress only on these scalar write requests:
// Set, SetIfAbsent, SetIfPresent, SetIfEqual, SetIfNotEqual, SetIfAbsentOrEqual, SetIfPresentAndNotEqual,
// SetWithHash, SetIfPresentAndHashEqual, SetIfPresentAndHashNotEqual, SetIfAbsentOrHashEqual, SetIfAbsentOrHashNotEqual.
// Specify IncludeTypes in CompressionMiddlewareProps if you wish to compress only a subset of these requests.
func (rh *CompressionMiddlewareRequestHandler) OnRequest(req interface{}) (interface{}, error) {
	// We still need to use a switch statement to be able to access the request objects
	// as the specific request types in order to access the Value field.
	switch r := req.(type) {
	case *momento.SetRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfAbsentRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfPresentRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfNotEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfNotEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfAbsentOrEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndNotEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfPresentAndNotEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetWithHashRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetWithHashRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndHashEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfPresentAndHashEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfPresentAndHashNotEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfPresentAndHashNotEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrHashEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfAbsentOrHashEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
		}
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrHashNotEqualRequest:
		rawData := []byte(fmt.Sprintf("%v", r.Value))
		compressed, err := rh.compressor.Compress(rawData, rh.GetLogger(), "SetIfAbsentOrHashNotEqualRequest")
		if err != nil {
			return nil, fmt.Errorf("failed to compress request: %v", err)
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
		decompressed, err := rh.compressor.Decompress(rawData, rh.GetLogger(), "GetHit")
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		return responses.NewGetHit(decompressed), nil
	case *responses.GetWithHashHit:
		rawData := r.ValueByte()
		decompressed, err := rh.compressor.Decompress(rawData, rh.GetLogger(), "GetWithHashHit")
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		return responses.NewGetWithHashHit(decompressed, r.HashByte()), nil
	default:
		rh.GetLogger().Info("No action for OnResponse type: %T", resp)
		return resp, nil
	}
}
