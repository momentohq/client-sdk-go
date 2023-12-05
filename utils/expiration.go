package utils

import (
	"math"
	"time"
)

type Expiration interface {
	DoesExpire() bool
}

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
