package utils

import (
	hdr "github.com/HdrHistogram/hdrhistogram-go"
	"time"
)

type PerfTestContext struct {
	StartTime          time.Time
	TotalItemSizeBytes int64
	AsyncGetLatencies  *hdr.Histogram
	AsyncSetLatencies  *hdr.Histogram
	SetBatchLatencies  *hdr.Histogram
	GetBatchLatencies  *hdr.Histogram
}

func InitiatePerfTestContext() *PerfTestContext {
	return &PerfTestContext{
		StartTime:          time.Now(),
		TotalItemSizeBytes: 0,
		AsyncGetLatencies:  hdr.New(1, 10000000, 3),
		AsyncSetLatencies:  hdr.New(1, 10000000, 3),
		SetBatchLatencies:  hdr.New(1, 10000000, 3),
		GetBatchLatencies:  hdr.New(1, 10000000, 3),
	}
}

type RequestType string

const (
	ASYNC_GETS RequestType = "ASYNC_GETS"
	ASYNC_SETS RequestType = "ASYNC_SETS"
	GET_BATCH  RequestType = "GET_BATCH"
	SET_BATCH  RequestType = "SET_BATCH"
)

type PerfTestOptions struct {
	RequestTimeoutSeconds time.Duration
}

type GetSetConfig struct {
	BatchSize     int
	ItemSizeBytes int
}

type PerfTestConfiguration struct {
	MinimumRunDurationSecondsForTests int
	Sets                              []GetSetConfig
	Gets                              []GetSetConfig
}
