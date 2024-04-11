package utils

import (
	"github.com/HdrHistogram/hdrhistogram-go"
	"time"
)

type PerfTestContext struct {
	StartTime          time.Time
	TotalItemSizeBytes int64
	AsyncGetLatencies  *hdrhistogram.Histogram
	AsyncSetLatencies  *hdrhistogram.Histogram
	SetBatchLatencies  *hdrhistogram.Histogram
	GetBatchLatencies  *hdrhistogram.Histogram
}

func InitiatePerfTestContext() *PerfTestContext {
	return &PerfTestContext{
		StartTime:          time.Now(),
		TotalItemSizeBytes: 0,
		AsyncGetLatencies:  hdrhistogram.New(1, 1000, 1),
		AsyncSetLatencies:  hdrhistogram.New(1, 1000, 1),
		SetBatchLatencies:  hdrhistogram.New(1, 1000, 1),
		GetBatchLatencies:  hdrhistogram.New(1, 1000, 1),
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
