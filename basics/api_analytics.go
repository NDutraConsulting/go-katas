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
		"edge 200 1500",
		"edge 500 50",
		"auth 200 200",
		"auth 500 100",
		"user 200 100",
		// Error cases
		"edge 200 1x0",
		"edge 500",
		"user 200 ",
		"user",
		"200 1500",
		"",
	}
}

type internalApiInfo struct {
	successCount      int
	avgLatency        int
	validLatencyCount int
	latencyError      []string
}

/**
* Parse all logs first (stores parsed structs)
* Larger memory footprint
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
* Parse + process in one pass (doesn’t store parsed logs)
* Smaller memory footprint
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
	LatencyErrors     []string
}

func extractJsonApiMap(apiMap map[string]internalApiInfo) string {

	jsonReadyApiMap := make(map[string]PublicAPIInfo, len(apiMap))
	for key, value := range apiMap {
		jsonReadyApiMap[key] = PublicAPIInfo{
			SuccessCount:      value.successCount,
			AvgLatency:        value.avgLatency,
			ValidLatencyCount: value.validLatencyCount,
			LatencyErrors:     value.latencyError,
		}
	}
	jsonApiMap, err := json.Marshal(jsonReadyApiMap)
	if err != nil {
		return "{}"
	}
	return string(jsonApiMap)
}

type ParsedLog struct {
	Service      string
	Status       string
	Latency      int
	LatencyError string
}

func parseLogLine(line string) ParsedLog {
	line = strings.TrimSpace(line)
	f := strings.Fields(line)

	var out ParsedLog
	if len(f) == 0 {
		return out
	}

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
	} else if len(f) == 2 {
		out.LatencyError = "latency missing"
	}

	return out
}
