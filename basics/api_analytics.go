package basics

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

func runAPI() {

	logHistory := requestHistory()
	logs := make([][]string, len(logHistory))

	for i, log := range logHistory {
		logs[i] = strings.Fields(log)
	}

	apiMap := make(map[string]ApiInfo, len(logs))
	apiLatency := map[string]int{}

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

	outMap := &apiMap
	derefrenceApiMap, _ := json.Marshal(outMap)

	fmt.Println("JSON: ", string(derefrenceApiMap))

	fmt.Println(apiMap)

}
