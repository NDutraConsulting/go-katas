package robotfarm

import (
	"testing"
)

func TestRottenOranges(t *testing.T) {
	orchard := [][]int{
		{2, 1, 1},
		{1, 1, 0},
		{0, 1, 1},
	}
	result := orangesRotting(orchard)
	if result != 4 {
		t.Errorf("orchard(%v) - result time: %d; want %d", orchard, result, 3)
	}
}

func TestRottenOrangeRobots(t *testing.T) {
	orchard := [][]int{
		{2, 1, 1},
		{1, 3, 1},
		{1, 1, 2},
	}
	result := orangesRotting(orchard)
	expect := -1
	if result != expect {
		t.Errorf("orchard(%v) - result time: %d; want %d", orchard, result, expect)
	}

	orchardB := [][]int{
		{2, 2, 3},
		{1, 1, 1},
		{1, 1, 1},
	}
	result = orangesRotting(orchardB)
	expect = 3
	if result != expect {
		t.Errorf("orchard(%v) - result time: %d; want %d", orchardB, result, expect)
	}

	orchardC := [][]int{
		{2, 2, 3, 1},
		{1, 1, 1, 1},
		{1, 3, 1, 1},
	}
	result = orangesRotting(orchardC)
	expect = -1
	if result != expect {
		t.Errorf("orchard(%v) - result time: %d; want %d", orchardC, result, expect)
	}

	orchardD := [][]int{
		{1, 2, 3, 1},
	}
	result = orangesRotting(orchardD)
	expect = -1
	if result != expect {
		t.Errorf("orchard(%v) - result time: %d; want %d", orchardD, result, expect)
	}
}
