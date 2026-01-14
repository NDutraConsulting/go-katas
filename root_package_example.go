package main

import (
	"fmt"
	"katas/basics"
	"math"
)

type box = basics.Box

var Roots = 333

func PackageTest() {
	fmt.Println("Im in the root directory.")
	fmt.Println("The sum is:", add(5, 3.49))

	fmt.Println("Boxes are:", makeBoxPile())
}

func add(a int, b float32) int {
	x := int(math.Round(float64(b)))
	return a + x
}

func makeBoxPile() []box {
	pile := []box{}
	for i := 0; i < 3; i++ {
		pile = append(pile, box{Length: 10 + i, Width: 5, Height: 3 + i})
	}
	return pile
}
