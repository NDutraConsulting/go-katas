package basics

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func requestHistory() []string {

	// api-name response latency
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
	count             int
	avgLatency        int
	validLatencyCount int
	latencyErrors     []string
}

/**
* Parse all logs first (stores parsed structs)
* Larger memory footprint
**/
func runHistoryAnalyticsA() (string, int64) {
	start := time.Now()
	logHistory := requestHistory()

	// Initialize the memory
	logs := make([]ParsedLog, len(logHistory))
	for i, log := range logHistory {
		logs[i] = parseLogLine(log)
	}

	apiMapSuccess := map[string]internalApiInfo{}
	apiMapFailure := map[string]internalApiInfo{}
	apiLatency := map[string]int{}

	for _, logArr := range logs {
		switch logArr.Status {
		case "200":
			setData(logArr, apiMapSuccess, apiLatency)
		case "500":
			setData(logArr, apiMapFailure, apiLatency)
		}
	}

	elapsed := time.Since(start)
	return extractJsonApiMap(apiMapSuccess), elapsed.Nanoseconds()
}

/**
* Parse + process in one pass (doesn’t store parsed logs)
* Smaller memory footprint
**/
func runHistoryAnalyticsB() (string, int64) {
	start := time.Now()

	logHistory := requestHistory()
	apiMapSuccess := map[string]internalApiInfo{}
	apiMapFailure := map[string]internalApiInfo{}
	apiLatency := map[string]int{}
	for _, log := range logHistory {

		logArr := parseLogLine(log)
		switch logArr.Status {
		case "200":
			setData(logArr, apiMapSuccess, apiLatency)
		case "500":
			setData(logArr, apiMapFailure, apiLatency)
		}
	}

	elapsed := time.Since(start)
	return extractJsonApiMap(apiMapSuccess), elapsed.Nanoseconds()
}

func setData(logArr ParsedLog, apiMap map[string]internalApiInfo, apiLatency map[string]int) {
	if logArr.Service == "" {
		return
	}

	key := logArr.Service
	entry, exists := apiMap[key]
	if !exists {
		entry = internalApiInfo{count: 1}
	} else {
		entry.count++
	}

	if logArr.LatencyError != "" {
		entry.latencyErrors = append(entry.latencyErrors, logArr.LatencyError)
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
	LatencyErrors     []string `json:"latency_errors,omitempty"`
}

func extractJsonApiMap(apiMap map[string]internalApiInfo) string {

	jsonReadyApiMap := make(map[string]PublicAPIInfo, len(apiMap))
	for key, value := range apiMap {
		jsonReadyApiMap[key] = PublicAPIInfo{
			SuccessCount:      value.count,
			AvgLatency:        value.avgLatency,
			ValidLatencyCount: value.validLatencyCount,
			LatencyErrors:     value.latencyErrors,
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
