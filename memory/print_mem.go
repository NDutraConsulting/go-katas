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

func PrintKBMem(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf(
		"%s | Alloc=%dKB HeapAlloc=%dKB Sys=%dKB NumGC=%d\n",
		label,
		m.Alloc/1024,
		m.HeapAlloc/1024,
		m.Sys/1024,
		m.NumGC,
	)
}
