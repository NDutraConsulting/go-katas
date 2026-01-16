package leetcode

import (
	"testing"
)

func TestMaxDepth(t *testing.T) {
	tests := []treeTest{
		{tree: BuildTree([]int{5}), expect: 1},
		{tree: BuildTree([]int{1, 2, 3, 4, 5}), expect: 3},
		{tree: BuildTree([]int{7, 22, 3, -4, 5, 9, 99, 100, 101}), expect: 4},
	}
	for _, tc := range tests {
		if got := maxDepth(tc.tree); got != tc.expect {
			t.Errorf("maxDepth(%v) = %d; want %d", tc.tree, got, tc.expect)
		}
	}
}
