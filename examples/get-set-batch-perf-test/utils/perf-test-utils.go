package utils

import (
	"fmt"
	hdr "github.com/HdrHistogram/hdrhistogram-go"
	"os"
	"strings"
)

func outputHistogramSummary(histogram *hdr.Histogram) string {
	return fmt.Sprintf(`
    count: %d
    min: %d
    p50: %d
    p90: %d
    p99: %d
  p99.9: %d
    max: %d
`, histogram.TotalCount(), histogram.Min(), histogram.ValueAtPercentile(50), histogram.ValueAtPercentile(90), histogram.ValueAtPercentile(99), histogram.ValueAtPercentile(99.9), histogram.Max())
}

func CalculateSummary(context *PerfTestContext, batchSize int, itemSizeBytes int, requestType RequestType) {
	var histogram *hdr.Histogram
	switch requestType {
	case ASYNC_SETS:
		histogram = context.AsyncSetLatencies
	case ASYNC_GETS:
		histogram = context.AsyncGetLatencies
	case SET_BATCH:
		histogram = context.SetBatchLatencies
	case GET_BATCH:
		histogram = context.GetBatchLatencies
	}

	summaryMessage := generateSummaryMessage(requestType, histogram, batchSize, itemSizeBytes, context.TotalItemSizeBytes)

	// print the summary message
	fmt.Println(summaryMessage)

	// Write the statistics to a CSV file
	writeStatsToCSV(requestType, batchSize, itemSizeBytes, histogram, context.TotalItemSizeBytes)
}

func generateSummaryMessage(requestType RequestType, histogram *hdr.Histogram, batchSize int, itemSizeBytes int, totalItemSizeBytes int64) string {
	summaryTitle := fmt.Sprintf("======= Summary of %s requests for batch size %d and item size %d bytes =======", requestType, batchSize, itemSizeBytes)
	histogramSummary := fmt.Sprintf("Cumulative latencies: %s", outputHistogramSummary(histogram))
	totalItemSize := fmt.Sprintf("Total item size in bytes: %d bytes", totalItemSizeBytes)
	separator := fmt.Sprintf("%s\n\n", strings.Repeat("=", 150))

	return fmt.Sprintf("%s\n%s\n%s\n%s\n", summaryTitle, histogramSummary, totalItemSize, separator)
}

func writeStatsToCSV(requestType RequestType, batchSize int, itemSize int, histogram *hdr.Histogram, totalItemSizeBytes int64) {
	filename := "perf_test_stats.csv"
	header := "requestType,BatchSize,itemSize,TotalCount,Min,p50,p90,p99,p99.9,Max,TotalItemSizeBytes\n"
	stats := fmt.Sprintf("%s,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n", requestType, batchSize, itemSize, histogram.TotalCount(), histogram.Min(), histogram.ValueAtPercentile(50), histogram.ValueAtPercentile(90), histogram.ValueAtPercentile(99), histogram.ValueAtPercentile(99.9), histogram.Max(), totalItemSizeBytes)

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// If the file doesn't exist, write the header
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}(file)
		if _, err := file.WriteString(header); err != nil {
			fmt.Printf("Error writing header to file: %v\n", err)
		}
	}

	// Append the statistics to the file
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(file)
	if _, err := file.WriteString(stats); err != nil {
		fmt.Printf("Error writing stats to file: %v\n", err)
	}
}
