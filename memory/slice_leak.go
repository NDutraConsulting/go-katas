package memory

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

func createSlice() []int {
	return make([]int, 10000)
}

func getValue(s []int) []int {
	res := s[:3]
	return res
}

func getValueSlice(s []int) []int {
	res := make([]int, 3)
	copy(res, s[:3])
	return res
}

func memCapsule() {

	fmt.Println("\n\n Use getValue() <  retains a reference on each execution of the for loop scope.")

	for i := 0; i < 10; i++ {
		s := createSlice()
		val := getValue(s)

		fmt.Printf("%p\n", &val)
		PrintMem("tick: " + strconv.Itoa(i))

		// Drop references will clean this up
		val = nil
		s = nil
		runtime.GC()

		time.Sleep(100 * time.Millisecond)
	}

	runtime.GC()
	debug.FreeOSMemory() // optional, more aggressive
	PrintMem("getValue() End after Garbag Collector call.\n")

	fmt.Println("\n\n Use > getValueSlice() < is recommended but does not work...")
	fmt.Println("We must set the variables in the loop to nil to free memory.")

	for i := 0; i < 10; i++ {
		s := createSlice()
		val := getValueSlice(s)
		fmt.Printf("%p\n", &val)
		PrintMem("tick: " + strconv.Itoa(i))
		time.Sleep(100 * time.Millisecond)
	}

	runtime.GC()
	debug.FreeOSMemory() // optional, more aggressive
	PrintMem("getValueSlice() End after Garbag Collector call.\n")
}

func RunSliceLeakDemo() {

	fmt.Println("This is a scenario that can occur when polling a PORT, endpoint or kafka channel.")
	PrintMem("start")
	memCapsule()

}
