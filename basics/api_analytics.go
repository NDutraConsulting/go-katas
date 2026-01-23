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
		"edge 500",
		"auth 500 100",
		"user 200 100",
		"edge 200 1500",
	}
}

type internalApiInfo struct {
	successCount      int
	avgLatency        int
	validLatencyCount int
	latencyError      []string
}

/**
* Convert the log history into a map of API information BEFORE processing
* This has a larger memory footprint
**/
func runHistoryAnaliticsA() (string, int64) {
	start := time.Now()
	logHistory := requestHistory()

	// Initialize the memory
	logs := make([]ParsedLog, len(logHistory))
	for i, log := range logHistory {
		logs[i] = parseLogLine(log)
	}

	apiMap := map[string]internalApiInfo{}
	apiLatency := map[string]int{}
	for _, logArr := range logs {

		if logArr.Status != "200" {
			// We only care about successful requests
			continue
		}

		setDataForSuccess(logArr, apiMap, apiLatency)
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

		logArr := parseLogLine(log)
		if logArr.Status != "200" {
			// We only care about successful requests
			continue
		}

		setDataForSuccess(logArr, apiMap, apiLatency)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	return extractJsonApiMap(apiMap), elapsed.Nanoseconds()
}
func setDataForSuccess(logArr ParsedLog, apiMap map[string]internalApiInfo, apiLatency map[string]int) {

	key := logArr.Service
	entry, exists := apiMap[key]
	if !exists {
		entry = internalApiInfo{successCount: 1}
	} else {
		entry.successCount++
	}

	if logArr.LatencyError != "" {
		entry.latencyError = append(entry.latencyError, logArr.LatencyError)
		apiMap[key] = entry
		// Do not continue processing latency
		return
	}

	apiLatency[key] += logArr.Latency
	entry.validLatencyCount++
	entry.avgLatency = apiLatency[key] / entry.validLatencyCount

	apiMap[key] = entry
}

type PublicAPIInfo struct {
	SuccessCount      int
	AvgLatency        int
	ValidLatencyCount int
	LatencyError      []string
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

type ParsedLog struct {
	Service      string
	Status       string
	Latency      int
	LatencyError string
}

func parseLogLine(line string) ParsedLog {
	f := strings.Fields(line) // handles extra spaces/tabs safely
	var out ParsedLog

	if len(f) >= 1 {
		out.Service = f[0]
	}
	if len(f) >= 2 {
		out.Status = f[1]
	}
	if len(f) >= 3 {
		l, err := strconv.Atoi(f[2])
		if err != nil {
			out.LatencyError = err.Error()
		} else {
			out.Latency = l
		}
	}
	return out
}
