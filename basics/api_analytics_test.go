package basics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestAPIAnalytics(t *testing.T) {

	responseA, elapsedA := runHistoryAnalyticsA()
	responseB, elapsedB := runHistoryAnalyticsB()

	if responseA != responseB {
		t.Errorf("Expected same output, got different outputs")
	}

	// Accounting for map ordering issues
	var mapA map[string]PublicAPIInfo
	var mapB map[string]PublicAPIInfo
	err := json.Unmarshal([]byte(responseA), &mapA)
	if err != nil {
		t.Fatalf("failed to unmarshal responseA: %v", err)
	}

	err = json.Unmarshal([]byte(responseB), &mapB)
	if err != nil {
		t.Fatalf("failed to unmarshal responseB: %v", err)
	}

	if !reflect.DeepEqual(mapA, mapB) {
		t.Errorf("Expected same output, got different outputs\nA=%v\nB=%v", mapA, mapB)
	}

	fmt.Println("responseA:", responseA)
	fmt.Println("responseB:", responseB)

	fmt.Println("------------ RUN 1 (A -> B)------------")
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("------------ RUN 2 (B -> A)------------")
	responseB, elapsedB = runHistoryAnalyticsB()
	responseA, elapsedA = runHistoryAnalyticsA()
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("------------ RUN 3 (B -> A)------------")
	responseB, elapsedB = runHistoryAnalyticsB()
	responseA, elapsedA = runHistoryAnalyticsA()
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("------------ RUN 4 (A -> B)------------")
	responseA, elapsedA = runHistoryAnalyticsA()
	responseB, elapsedB = runHistoryAnalyticsB()
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("\n------------ RUN Averages ------------")
	runs := int64(100000)
	fmt.Println("Run each function: ", runs, " times and take the average.")

	aAvg := runTestBlock(runHistoryAnalyticsA, runs)
	bAvg := runTestBlock(runHistoryAnalyticsB, runs)

	fmt.Println("A Avg: ", aAvg, "ns --- B Avg: ", bAvg, "ns")
}

func runTestBlock(f func() (string, int64), runs int64) int64 {
	ch := make(chan int64)

	go func(f func() (string, int64), runs int64) {
		total := int64(0)
		for range runs {
			total += timeTest(f)
		}
		ch <- total / runs
	}(f, runs)

	return <-ch
}

func timeTest(f func() (string, int64)) int64 {

	_, time := f()

	return time
}
