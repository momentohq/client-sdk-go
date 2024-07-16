package utils

import (
	"math"
	"time"
)

type Expiration interface {
	DoesExpire() bool
}

// ExpiresIn is used in AuthClient requests to specify the expiration time of a disposable token or API key.
type ExpiresIn struct {
	validForSeconds int64
}

func (e ExpiresIn) DoesExpire() bool {
	return e.validForSeconds != math.MaxInt64
}

func (e ExpiresIn) Seconds() int64 {
	return e.validForSeconds
}

func ExpiresInNever() ExpiresIn {
	return ExpiresIn{validForSeconds: math.MaxInt64}
}

func ExpiresInSeconds(seconds int64) ExpiresIn {
	return ExpiresIn{validForSeconds: seconds}
}

func ExpiresInMinutes(minutes int64) ExpiresIn {
	return ExpiresIn{validForSeconds: minutes * 60}
}

func ExpiresInHours(hours int64) ExpiresIn {
	return ExpiresIn{validForSeconds: hours * 60 * 60}
}

func ExpiresInDays(days int64) ExpiresIn {
	return ExpiresIn{validForSeconds: days * 60 * 60 * 24}
}

func ExpiresAtEpoch(expiresBy int64) ExpiresIn {
	var now = time.Now().Unix()
	return ExpiresIn{validForSeconds: expiresBy - now}
}

// ExpiresAt is used in AuthClient responses to specify the expiration time of a disposable token or API key.
type ExpiresAt struct {
	validUntil int64
}

func (e ExpiresAt) DoesExpire() bool {
	return e.validUntil != math.MaxInt64
}

func (e ExpiresAt) Epoch() int64 {
	return e.validUntil
}

// ExpiresAtFromEpoch constructs an ExpiresAt with the specified epoch timestamp,
// but if timestamp is undefined, the epoch timestamp will be set to math.MaxInt64.
func ExpiresAtFromEpoch(epoch *int64) ExpiresAt {
	if epoch != nil {
		return ExpiresAt{validUntil: *epoch}
	}
	return ExpiresAt{validUntil: math.MaxInt64}
}
