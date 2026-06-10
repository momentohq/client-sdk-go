package momentoerrors

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMomentoSvcErrorUnwrapExposesCause(t *testing.T) {
	t.Parallel()

	cause := context.Canceled
	err := NewMomentoSvcErr(CanceledError, "subscribe context canceled", cause)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("errors.Is(err, context.Canceled) = false, want true (err: %v)", err)
	}

	var svcErr MomentoSvcErr
	if !errors.As(err, &svcErr) || svcErr.Code() != CanceledError {
		t.Fatalf("errors.As did not recover the MomentoSvcErr (err: %v)", err)
	}
}

func TestMomentoSvcErrorUnwrapNilCause(t *testing.T) {
	t.Parallel()

	err := NewMomentoSvcErr(TimeoutError, "no cause", nil)
	if errors.Unwrap(err) != nil {
		t.Fatalf("Unwrap() = %v, want nil", errors.Unwrap(err))
	}
}

// TestConvertSvcErrUnwrapsToGrpcStatus pins the behavior change Unwrap
// introduces in ConvertSvcErr: a MomentoSvcErr wrapping a gRPC status error
// now re-converts by the inner status code instead of falling back to
// ClientSdkError.
func TestConvertSvcErrUnwrapsToGrpcStatus(t *testing.T) {
	t.Parallel()

	wrapped := NewMomentoSvcErr(ServerUnavailableError, "stream disconnected", status.Error(codes.Unavailable, "unavailable"))
	if got := ConvertSvcErr(wrapped).Code(); got != ServerUnavailableError {
		t.Fatalf("ConvertSvcErr(wrapped grpc status).Code() = %s, want %s", got, ServerUnavailableError)
	}

	plain := NewMomentoSvcErr(CanceledError, "subscription closed", nil)
	if got := ConvertSvcErr(plain).Code(); got != ClientSdkError {
		t.Fatalf("ConvertSvcErr(no-status MomentoSvcErr).Code() = %s, want %s (documented fallback)", got, ClientSdkError)
	}
}
