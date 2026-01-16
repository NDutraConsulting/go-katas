package memory

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v KB\n", bToKb(m.Alloc))
}

func bToKb(b uint64) uint64 {
	return b / 1024
}

type options struct {
	clearMap bool
}

func MapLeak(op options) {
	n := 200000
	m := make(map[int][128]string)

	printAlloc()

	for i := 0; i < n; i++ {
		m[i] = [128]string{}
	}

	printAlloc()

	for i := range n {
		delete(m, i)
	}

	runtime.GC()

	printAlloc()

	// Keep Alive up to this point
	runtime.KeepAlive(m)

	if op.clearMap {
		m = nil
		runtime.GC()
		printAlloc()
	}

	m = nil
}

func RunMapLeakTest() {
	runA()
	scopeShiftTrigger()
	runB()
}

func runA() {
	fmt.Println("Start of runA")
	fmt.Println("> No runtime.GC() call but we do set m=nil and wait for GC...<")
	PrintKBMem(">")

	fmt.Println("\nDoes NOT Clear Map with m=nil even after waiting...")
	fmt.Println(" - WE must set m=nil; and then call runtime.GC()")
	MapLeak(options{clearMap: false})
	l := make(map[int]string)

	seconds := 5
	fmt.Printf("\nwaiting %v seconds...", seconds)
	for i := 0; i < seconds; i++ {
		fmt.Print(".")
		time.Sleep(1000 * time.Millisecond)
		l[i] = "test" + strconv.Itoa(i)
	}
	fmt.Println(l)
	printAlloc()
	fmt.Println("End of runA")
}

func runB() {
	fmt.Println("\n\nStart of runB")
	fmt.Println("> Use runtime.GC() call <")
	fmt.Println("Clears Map with m=nil & runtime.GC()")
	MapLeak(options{clearMap: true})
	printAlloc()
	PrintKBMem(">")
	fmt.Println("End of runB")

}

func scopeShiftTrigger() {
	fmt.Println("\n\n Run > scopeShiftTrigger()")
	x := "Scoped String"
	fmt.Println(x)
	PrintKBMem(">")
	fmt.Println("m := make(map[int][128]string) is still being held in the heap.")
}
