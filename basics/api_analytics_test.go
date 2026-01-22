package basics

import (
	"fmt"
	"testing"
)

func TestAPIAnalytics(t *testing.T) {

	fmt.Println("------------ RUN 1 ------------")
	responseA, elapsedA := runHistoryAnaliticsA()
	responseB, elapsedB := runHistoryAnaliticsB()

	fmt.Println("A time: ", elapsedA, "response:", responseA)
	fmt.Println("B time: ", elapsedB, "response:", responseB)
	if responseA != responseB {
		t.Errorf("Expected same output, got different outputs")
	}

	fmt.Println("------------ RUN 2 ------------")
	responseB, elapsedB = runHistoryAnaliticsB()
	responseA, elapsedA = runHistoryAnaliticsA()
	fmt.Println("A time: ", elapsedA, "response:", responseA)
	fmt.Println("B time: ", elapsedB, "response:", responseB)

	fmt.Println("------------ RUN 3 ------------")
	responseB, elapsedB = runHistoryAnaliticsB()
	responseA, elapsedA = runHistoryAnaliticsA()
	fmt.Println("A time: ", elapsedA, "response:", responseA)
	fmt.Println("B time: ", elapsedB, "response:", responseB)

	fmt.Println("------------ RUN 4 ------------")
	responseA, elapsedA = runHistoryAnaliticsA()
	responseB, elapsedB = runHistoryAnaliticsB()
	fmt.Println("A time: ", elapsedA, "response:", responseA)
	fmt.Println("B time: ", elapsedB, "response:", responseB)

	fmt.Println("\n------------ RUN Averages ------------")
	runs := int64(100000)
	fmt.Println("Run each function: ", runs, " times and take the average.")

	aAvg := runTestBlock(runHistoryAnaliticsA, runs)
	bAvg := runTestBlock(runHistoryAnaliticsB, runs)

	bAvg += runTestBlock(runHistoryAnaliticsB, runs)
	aAvg += runTestBlock(runHistoryAnaliticsA, runs)

	bAvg += runTestBlock(runHistoryAnaliticsB, runs)
	aAvg += runTestBlock(runHistoryAnaliticsA, runs)

	aAvg = aAvg / 3
	bAvg = bAvg / 3

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
