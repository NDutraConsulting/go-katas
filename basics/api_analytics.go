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
		"edge 200 100",
		"edge 500 50",
		"auth 200 200",
		"edge 500 50",
		"auth 500 100",
		"user 200 100",
	}
}

type ApiInfo struct {
	SuccessCount int
	AvgLatency   int
}

func runHistoryAnaliticsA() (string, int64) {

	start := time.Now()
	logHistory := requestHistory()

	// We know the number of elements we need so lets initialize the memory
	logs := make([][]string, len(logHistory))

	for i, log := range logHistory {
		logs[i] = strings.Fields(log)
	}

	// This might have only 1 entry
	apiMap := map[string]ApiInfo{}

	// We know the number of elements we need so lets initialize the memory
	apiLatency := make(map[string]int, len(logs))

	for _, logArr := range logs {

		key := logArr[0]
		latency, _ := strconv.Atoi(logArr[2])
		apiLatency[key] += latency

		_, keyExists := apiMap[key]
		if keyExists {

			entry := apiMap[key]
			entry.SuccessCount++

			avgLatency := apiLatency[key] / entry.SuccessCount
			entry.AvgLatency = avgLatency

			apiMap[key] = entry

			continue
		}

		apiMap[key] = ApiInfo{
			SuccessCount: 1,
			AvgLatency:   latency,
		}

	}

	t := time.Now()
	elapsed := t.Sub(start)

	//fmt.Println("\n------------- runHistoryAnaliticsFast() -------------")
	//fmt.Println("Good memory management results -> ", elapsed)

	outMap := &apiMap
	jsonApiMap, _ := json.Marshal(outMap)

	return string(jsonApiMap), elapsed.Nanoseconds()

}

func runHistoryAnaliticsB() (string, int64) {
	start := time.Now()

	logHistory := requestHistory()

	// This might have only 1 entry
	apiMap := map[string]ApiInfo{}

	// We know the number of elements we need so lets initialize the memory
	apiLatency := make(map[string]int, len(logHistory))

	for _, log := range logHistory {

		logArr := strings.Fields(log)

		key := logArr[0]
		latency, _ := strconv.Atoi(logArr[2])
		apiLatency[key] += latency

		_, keyExists := apiMap[key]
		if keyExists {

			entry := apiMap[key]
			entry.SuccessCount++

			avgLatency := apiLatency[key] / entry.SuccessCount
			entry.AvgLatency = avgLatency

			apiMap[key] = entry

			continue
		}

		apiMap[key] = ApiInfo{
			SuccessCount: 1,
			AvgLatency:   latency,
		}

	}

	t := time.Now()
	elapsed := t.Sub(start)
	//fmt.Println("\n------------- runHistoryAnaliticsSlow() -------------")
	//fmt.Println("Bad memory management results -> ", elapsed)

	outMap := &apiMap
	jsonApiMap, _ := json.Marshal(outMap)

	return string(jsonApiMap), elapsed.Nanoseconds()

}
