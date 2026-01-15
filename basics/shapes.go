package basics

import (
	"fmt"
	"path/filepath"
	"runtime"
)

var f = fmt.Println

type shape struct {
	id    int8
	label string
}

func (s *shape) setLabel(label string) {
	s.label = label
}

type rect struct {
	shape
	length uint32
	width  uint32
	area   uint64
}

// This will not set r.area because it copies r.
func (r rect) getArea() uint64 {
	printFuncName()
	fmt.Println("WARNING: This will not set r.area because it copies r.")
	fmt.Printf("Address of r = %p\n", &r)

	if r.area != 0 {
		fmt.Println("CACHE reached!")
		return r.area
	}

	fmt.Println("Area is 0 we must compute i and update the memory")

	r.area = uint64(r.length) * uint64(r.width)

	return r.area
}

// Dereference so that we can mutate the rect
func (r *rect) getCachedArea() uint64 {
	printFuncName()
	fmt.Printf("Address of r = %p\n", r)

	fmt.Println("get")
	if r.area != 0 {
		fmt.Println("CACHE reached!")
		return r.area
	}

	fmt.Println("Area is 0 we must compute i and update the memory")

	r.area = uint64(r.length) * uint64(r.width)

	return r.area
}

func RunShapeDemo() {

	r := rect{
		shape:  shape{id: 1, label: "x"},
		length: 10,
		width:  5,
		area:   0,
	}

	fmt.Printf("Address of r = %p\n", &r)

	r.shape.setLabel("rect1")

	fmt.Println(r.shape.label)

	fmt.Println(r.getArea())

	fmt.Println(r.getArea())

	fmt.Println(r.getCachedArea())
	fmt.Println(r.getCachedArea())

}

func printFuncName() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()

	rv := filepath.Base(funcName)

	fmt.Printf("\nCurrently in function: %s\n", rv)
}
