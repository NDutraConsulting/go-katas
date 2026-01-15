package leetcode

import (
	"fmt"
	"testing"
)

type test struct {
	tree   *TreeNode
	expect int
}

func TestMaxLevelSum(t *testing.T) {

	var tests = []test{
		{tree: &TreeNode{
			Val:   1,
			Left:  &TreeNode{Val: -7, Right: nil, Left: nil},
			Right: &TreeNode{Val: 0, Right: nil, Left: nil}},
			expect: 1},
		{tree: BuildTree([]int{1, 7, 0, 7, -8}), expect: 2},
		{tree: BuildTree([]int{5}), expect: 1},
	}

	for idx := range tests {
		fmt.Println()
		result := maxLevelSum(tests[idx].tree)
		if result != tests[idx].expect {
			t.Errorf("Test idx: %v - Expected %v, got %v", idx, tests[idx], result)
		}
	}

}
