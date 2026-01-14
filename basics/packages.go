package basics

import (
	"fmt"
	"math/rand"
)

type Box struct {
	Length int
	Width  int
	Height int
}

func Package(boxes []Box) {
	fmt.Println(boxes)
}

func test() {
	fmt.Println("My favorite number is ", rand.Intn(10))
}
