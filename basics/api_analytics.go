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
		"edge 201 1500",
		"edge 202 1500",
		"edge 500 50",
		"auth 200 200",
		"auth 500 100",
		"user 200 100",

		// Edge cases
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
	sumLatency        int
	validLatencyCount int
	latencyErrors     []string
}

/**
* Parse all logs first (stores parsed structs)
* Larger memory footprint
**/
func runHistoryAnalyticsA() (string, string, int64) {
	start := time.Now()
	logHistory := requestHistory()

	// Initialize the memory
	logs := make([]ParsedLog, len(logHistory))
	for i, log := range logHistory {
		logs[i] = parseLogLine(log)
	}

	apiMapSuccess := map[string]internalApiInfo{}
	apiMapFailure := map[string]internalApiInfo{}

	for _, logArr := range logs {
		if !logArr.StatusValid {
			continue
		}
		if is2xx(logArr.Status) {
			setData(logArr, apiMapSuccess)
		} else if is5xx(logArr.Status) {
			setData(logArr, apiMapFailure)
		}
	}

	elapsed := time.Since(start)
	return extractJsonApiMap(apiMapSuccess, "success"), extractJsonApiMap(apiMapFailure, "failure"), elapsed.Nanoseconds()
}

/**
* Parse + process in one pass (doesn’t store parsed logs)
* Smaller memory footprint
**/
func runHistoryAnalyticsB() (string, string, int64) {
	start := time.Now()

	logHistory := requestHistory()
	apiMapSuccess := map[string]internalApiInfo{}

	apiMapFailure := map[string]internalApiInfo{}

	for _, log := range logHistory {
		logArr := parseLogLine(log)
		if !logArr.StatusValid {
			continue
		}

		if is2xx(logArr.Status) {
			setData(logArr, apiMapSuccess)
		} else if is5xx(logArr.Status) {
			setData(logArr, apiMapFailure)
		}
	}

	elapsed := time.Since(start)
	return extractJsonApiMap(apiMapSuccess, "success"), extractJsonApiMap(apiMapFailure, "failure"), elapsed.Nanoseconds()
}

func is5xx(s int) bool { return s >= 500 && s < 600 }
func is2xx(s int) bool { return s >= 200 && s < 300 }

func setData(logArr ParsedLog, apiMap map[string]internalApiInfo) {
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

	entry.sumLatency += logArr.Latency
	entry.validLatencyCount++
	entry.avgLatency = entry.sumLatency / entry.validLatencyCount

	apiMap[key] = entry
}

type PublicAPIInfo struct {
	Count             int
	AvgLatency        int
	ValidLatencyCount int
	LatencyErrors     []string `json:"latency_errors,omitempty"`
}

func extractJsonApiMap(apiMap map[string]internalApiInfo, responseType string) string {

	jsonReadyApiMap := make(map[string]PublicAPIInfo, len(apiMap))
	for key, value := range apiMap {
		jsonReadyApiMap[key] = PublicAPIInfo{
			Count:             value.count,
			AvgLatency:        value.avgLatency,
			ValidLatencyCount: value.validLatencyCount,
			LatencyErrors:     value.latencyErrors,
		}
	}

	rv := ResponseJSON{
		ResponseType: responseType,
		Services:     jsonReadyApiMap,
	}

	jsonApiMap, err := json.Marshal(rv)
	if err != nil {
		return "{}"
	}
	return string(jsonApiMap)
}

type ResponseJSON struct {
	ResponseType string                   `json:"response_type"`
	Services     map[string]PublicAPIInfo `json:"services"`
}

type ParsedLog struct {
	Service      string
	Status       int
	StatusValid  bool
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
		if s, err := strconv.Atoi(f[1]); err == nil {
			out.Status = s
			out.StatusValid = true
		}
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
