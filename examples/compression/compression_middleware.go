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

func (rh *compressionMiddlewareRequestHandler) OnRequest(req interface{}) (interface{}, error) {
	// Compress on writes
	switch r := req.(type) {
	case *momento.SetRequest:
		rh.GetLogger().Info("Setting key: %s, value: %s", r.Key, r.Value)
		rawData := r.Value
		compressed := rh.encoder.EncodeAll([]byte(fmt.Sprintf("%v", rawData)), nil)
		rh.GetLogger().Info("Compressed request %T: %d bytes -> %d bytes", req, len(fmt.Sprintf("%v", rawData)), len(compressed))
		return &momento.SetRequest{
			CacheName: r.CacheName,
			Key:       r.Key,
			Value:     momento.String(compressed),
		}, nil
	case *momento.SetIfAbsentOrHashEqualRequest:
		rh.GetLogger().Info("Setting key: %s, value: %s", r.Key, r.Value)
		rawData := r.Value
		compressed := rh.encoder.EncodeAll([]byte(fmt.Sprintf("%v", rawData)), nil)
		rh.GetLogger().Info("Compressed request %T: %d bytes -> %d bytes", req, len(fmt.Sprintf("%v", rawData)), len(compressed))
		return &momento.SetIfAbsentOrHashEqualRequest{
			CacheName: r.CacheName,
			Key:       r.Key,
			Value:     momento.String(compressed),
			HashEqual: r.HashEqual,
		}, nil
	default:
		rh.GetLogger().Info("No action for OnRequest type: %T", req)
		return req, nil
	}
}

func (rh *compressionMiddlewareRequestHandler) OnResponse(resp interface{}) (interface{}, error) {
	// Decompress on reads
	switch r := resp.(type) {
	case *responses.GetHit:
		rawData := r.ValueByte()
		decompressed, err := rh.decoder.DecodeAll(rawData, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		rh.GetLogger().Info(
			"Decompressed response %T: %d bytes -> %d bytes",
			resp, len(rawData), len(decompressed),
		)
		newGetResponse := responses.NewGetHit(decompressed)
		return newGetResponse, nil
	case *responses.GetWithHashHit:
		rawData := r.ValueByte()
		decompressed, err := rh.decoder.DecodeAll(rawData, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %v", err)
		}
		rh.GetLogger().Info(
			"Decompressed response %T: %d bytes -> %d bytes",
			resp, len(rawData), len(decompressed),
		)
		newGetWithHashResponse := responses.NewGetWithHashHit(decompressed, r.HashByte())
		return newGetWithHashResponse, nil
	default:
		rh.GetLogger().Info("No action for OnResponse type: %T", resp)
		return resp, nil
	}
}
