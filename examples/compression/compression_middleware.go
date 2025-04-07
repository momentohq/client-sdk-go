package main

import (
	"fmt"

	"github.com/klauspost/compress/zstd"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

type compressionMiddleware struct {
	middleware.Middleware
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func (mw *compressionMiddleware) GetRequestHandler(baseHandler middleware.RequestHandler) (middleware.RequestHandler, error) {
	return NewCompressionMiddlewareRequestHandler(baseHandler, mw.encoder, mw.decoder), nil
}

func NewCompressionMiddleware(props middleware.Props) middleware.Middleware {
	// We use the default compression level for the encoder, but this could be made configurable.
	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
	decoder, _ := zstd.NewReader(nil)
	mw := middleware.NewMiddleware(props)
	return &compressionMiddleware{mw, encoder, decoder}
}

type compressionMiddlewareRequestHandler struct {
	middleware.RequestHandler
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewCompressionMiddlewareRequestHandler(rh middleware.RequestHandler, encoder *zstd.Encoder, decoder *zstd.Decoder) middleware.RequestHandler {
	return &compressionMiddlewareRequestHandler{rh, encoder, decoder}
}

func (rh *compressionMiddlewareRequestHandler) compress(requestType string, rawData []byte) []byte {
	compressed := rh.encoder.EncodeAll(rawData, nil)
	rh.GetLogger().Info("Compressed request %s: %d bytes -> %d bytes", requestType, len(rawData), len(compressed))
	return compressed
}

func (rh *compressionMiddlewareRequestHandler) decompress(responseType string, rawData []byte) ([]byte, error) {
	decompressed, err := rh.decoder.DecodeAll(rawData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress response: %v", err)
	}
	rh.GetLogger().Info("Decompressed response %s: %d bytes -> %d bytes", responseType, len(rawData), len(decompressed))
	return decompressed, nil
}

// Compress on write requests
func (rh *compressionMiddlewareRequestHandler) OnRequest(req interface{}) (interface{}, error) {
	// We still need to use a switch statement to be able to access the request objects
	// as the specific request types in order to access the Value field.
	switch r := req.(type) {
	case *momento.SetRequest:
		compressed := rh.compress("SetRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	case *momento.SetIfAbsentOrHashEqualRequest:
		compressed := rh.compress("SetIfAbsentOrHashEqualRequest", []byte(fmt.Sprintf("%v", r.Value)))
		r.Value = momento.Bytes(compressed)
		return r, nil
	default:
		rh.GetLogger().Info("No action for OnRequest type: %T", req)
		return req, nil
	}
}

// Decompress on read responses
func (rh *compressionMiddlewareRequestHandler) OnResponse(resp interface{}) (interface{}, error) {
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
