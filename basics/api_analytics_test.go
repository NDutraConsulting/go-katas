package basics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestAPIAnalytics(t *testing.T) {

	successResponseA, failureResponseA, elapsedA := runHistoryAnalyticsA()
	successResponseB, failureResponseB, elapsedB := runHistoryAnalyticsB()

	if successResponseA != successResponseB {
		t.Errorf("Expected same output, got different outputs")
	}

	if failureResponseA != failureResponseB {
		t.Errorf("Expected same output, got different outputs")
	}

	// Accounting for map ordering issues
	var mapA map[string]PublicAPIInfo
	var mapB map[string]PublicAPIInfo
	err := json.Unmarshal([]byte(successResponseA), &mapA)
	if err != nil {
		t.Fatalf("failed to unmarshal successResponseA: %v", err)
	}

	err = json.Unmarshal([]byte(successResponseB), &mapB)
	if err != nil {
		t.Fatalf("failed to unmarshal successResponseB: %v", err)
	}

	if !reflect.DeepEqual(mapA, mapB) {
		t.Errorf("Expected same output, got different outputs\nA=%v\nB=%v", mapA, mapB)
	}

	fmt.Println("============================ Failure Analytics ============================")

	fmt.Println("failureResponseA:", failureResponseA)
	fmt.Println("failureResponseB:", failureResponseB)

	fmt.Println("\n============================ Success Analytics ============================")

	fmt.Println("successResponseA:", successResponseA)
	fmt.Println("successResponseB:", successResponseB)

	fmt.Println("\n------------ RUN 1 (A -> B)------------")
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("------------ RUN 2 (B -> A)------------")
	successResponseB, failureResponseB, elapsedB = runHistoryAnalyticsB()
	successResponseA, failureResponseA, elapsedA = runHistoryAnalyticsA()
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("------------ RUN 3 (B -> A)------------")
	successResponseB, failureResponseB, elapsedB = runHistoryAnalyticsB()
	successResponseA, failureResponseA, elapsedA = runHistoryAnalyticsA()
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("------------ RUN 4 (A -> B)------------")
	successResponseA, failureResponseA, elapsedA = runHistoryAnalyticsA()
	successResponseB, failureResponseB, elapsedB = runHistoryAnalyticsB()
	fmt.Println("A time: ", elapsedA)
	fmt.Println("B time: ", elapsedB)

	fmt.Println("\n------------ RUN Averages ------------")
	runs := int64(100000)
	fmt.Println("Run each function: ", runs, " times and take the average.")

	aAvg := runTestBlock(runHistoryAnalyticsA, runs)
	bAvg := runTestBlock(runHistoryAnalyticsB, runs)

	fmt.Println("A Avg: ", aAvg, "ns --- B Avg: ", bAvg, "ns")
}

func runTestBlock(f func() (string, string, int64), runs int64) int64 {
	ch := make(chan int64)

	go func(f func() (string, string, int64), runs int64) {
		total := int64(0)
		for range runs {
			total += timeTest(f)
		}
		ch <- total / runs
	}(f, runs)

	return <-ch
}

func timeTest(f func() (string, string, int64)) int64 {

	_, _, time := f()

	return time
}
