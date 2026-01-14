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
	var s = createSlice()

	fmt.Println("\n\n - getValue()")

	for i := 0; i < 15; i++ {
		s2 := createSlice()
		//fmt.Println(s[i])
		val := getValueSlice(s2)

		fmt.Printf("%p\n", &val)
		PrintMem("tick: " + strconv.Itoa(i))

		// drop references
		val = nil
		s2 = nil
		runtime.GC()

		time.Sleep(100 * time.Millisecond)
	}

	runtime.GC()

	debug.FreeOSMemory() // optional, more aggressive
	time.Sleep(5000 * time.Millisecond)

	fmt.Println("\n\n - getValueSlice()")

	for i := 0; i < 10; i++ {
		//fmt.Println(s[i])
		val := getValueSlice(s)
		fmt.Printf("%p\n", &val)
		PrintMem("tick: " + strconv.Itoa(i))
		time.Sleep(100 * time.Millisecond)
	}
}

func RunSliceLeakDemo() {

	PrintMem("start")
	memCapsule()
	time.Sleep(5000 * time.Millisecond)
	PrintMem("\nClose")

}
