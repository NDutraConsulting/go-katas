package memory

import (
	"fmt"
	"runtime"
)

func PrintMem(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf(
		"%s | Alloc=%dB HeapAlloc=%dB Sys=%dB NumGC=%d\n",
		label,
		m.Alloc,
		m.HeapAlloc,
		m.Sys,
		m.NumGC,
	)
}
