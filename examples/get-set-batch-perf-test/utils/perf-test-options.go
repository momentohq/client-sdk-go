package utils

import (
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/loov/hrtime"
)

type PerfTestContext struct {
	StartTime             time.Duration
	TotalItemSizeBytes    int64
	TotalNumberOfRequests int64
	AsyncGetLatencies     *hdrhistogram.Histogram
	AsyncSetLatencies     *hdrhistogram.Histogram
	SetBatchLatencies     *hdrhistogram.Histogram
	GetBatchLatencies     *hdrhistogram.Histogram
}

func InitiatePerfTestContext() *PerfTestContext {
	return &PerfTestContext{
		StartTime:             hrtime.Now(),
		TotalItemSizeBytes:    0,
		TotalNumberOfRequests: 0,
		AsyncGetLatencies:     hdrhistogram.New(1, 10000000000, 3),
		AsyncSetLatencies:     hdrhistogram.New(1, 10000000000, 3),
		SetBatchLatencies:     hdrhistogram.New(1, 10000000000, 3),
		GetBatchLatencies:     hdrhistogram.New(1, 10000000000, 3),
	}
}

type RequestType string

const (
	AsyncGets RequestType = "ASYNC_GETS"
	AsyncSets RequestType = "ASYNC_SETS"
	GetBatch  RequestType = "GET_BATCH"
	SetBatch  RequestType = "SET_BATCH"
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
