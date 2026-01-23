package basics

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func requestHistory() []string {

	// api-name responese latency
	return []string{
		"edge 200 1500",
		"edge 200 1x0",
		"edge 500 50",
		"auth 200 200",
		"edge 500 50",
		"auth 500 100",
		"user 200 100",
		"edge 200 1500",
	}
}

type internalApiInfo struct {
	successCount      int
	avgLatency        int
	validLatencyCount int
	latencyError      string
}

type PublicAPIInfo struct {
	SuccessCount      int
	AvgLatency        int
	ValidLatencyCount int
	LatencyError      string
}

/**
* Convert the log history into a map of API information BEFORE processing
* This has a larger memory footprint
**/
func runHistoryAnaliticsA() (string, int64) {

	start := time.Now()
	logHistory := requestHistory()

	// Initialize the memory
	logs := make([][]string, len(logHistory))
	for i, log := range logHistory {
		logs[i] = strings.Fields(log)
	}

	apiMap := map[string]internalApiInfo{}
	apiLatency := map[string]int{}
	for _, logArr := range logs {

		if logArr[1] != "200" {
			// We only care about successful requests
			continue
		}
		setData(logArr, apiMap, apiLatency)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	return extractJsonApiMap(apiMap), elapsed.Nanoseconds()
}

/**
* Convert the log history into a map of API information WHILE processing
* This reuses the same memory for the log string array
**/
func runHistoryAnaliticsB() (string, int64) {
	start := time.Now()

	logHistory := requestHistory()
	apiMap := map[string]internalApiInfo{}
	apiLatency := map[string]int{}
	for _, log := range logHistory {

		logArr := strings.Fields(log)
		if logArr[1] != "200" {
			// We only care about successful requests
			continue
		}
		setData(logArr, apiMap, apiLatency)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	return extractJsonApiMap(apiMap), elapsed.Nanoseconds()
}

func setData(logArr []string, apiMap map[string]internalApiInfo, apiLatency map[string]int) {

	key := logArr[0]
	latency, latencyErr := strconv.Atoi(logArr[2])

	_, keyExists := apiMap[key]
	if keyExists {

		entry := apiMap[key]
		entry.successCount++

		if latencyErr != nil {
			entry.latencyError = latencyErr.Error()
		} else {

			apiLatency[key] += latency
			entry.validLatencyCount++
			avgLatency := apiLatency[key] / entry.validLatencyCount
			entry.avgLatency = avgLatency

		}

		apiMap[key] = entry

		return
	}

	apiMap[key] = internalApiInfo{
		successCount:      1,
		avgLatency:        latency,
		validLatencyCount: initLatencyCount(latencyErr),
		latencyError:      getLatencyErrorVal(latencyErr),
	}
}

func getLatencyErrorVal(latencyErr error) string {
	if latencyErr != nil {
		return latencyErr.Error()
	}
	return ""
}

func initLatencyCount(latencyErr error) int {
	if latencyErr != nil {
		return 0
	}
	return 1
}

func extractJsonApiMap(apiMap map[string]internalApiInfo) string {

	jsonReadyApiMap := make(map[string]PublicAPIInfo, len(apiMap))
	for key, value := range apiMap {
		jsonReadyApiMap[key] = PublicAPIInfo{
			SuccessCount:      value.successCount,
			AvgLatency:        value.avgLatency,
			ValidLatencyCount: value.validLatencyCount,
			LatencyError:      value.latencyError,
		}
	}
	jsonApiMap, _ := json.Marshal(jsonReadyApiMap)
	return string(jsonApiMap)
}
