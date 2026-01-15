package leetcode

import (
	"fmt"
	"math"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func BuildTree(arr []int) *TreeNode {
	if len(arr) == 0 {
		return nil
	}
	root := &TreeNode{Val: arr[0], Right: nil, Left: nil}
	queue := []*TreeNode{root}
	i := 1
	for i < len(arr) && queue != nil {
		node := queue[0]
		queue = queue[1:]

		if i < len(arr) {
			node.Left = &TreeNode{Val: arr[i]}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(arr) {
			node.Right = &TreeNode{Val: arr[i]}
			queue = append(queue, node.Right)
		}
		i++
	}
	printTreeStructure("ROOT\n", root)
	return root
}

func printTreeStructure(label string, root *TreeNode) {

	if root == nil {
		return
	}

	printTreeStructure("Left", root.Left)

	printTreeStructure("Right", root.Right)
	fmt.Printf("Val :%d - Node %s\n", root.Val, label)

}

/**
*	Given the root of a binary tree, the level of its root is 1,
*   the level of its children is 1 and so on...
*	Return the level with the maximum sum of all nodes in that level.
 */
func maxLevelSum(root *TreeNode) int {

	q := []*TreeNode{root}
	maxSum := math.MinInt
	ans := 1
	level := 1

	// Look at all the nodes in the Queue
	for len(q) > 0 {

		size := len(q)
		sum := 0

		// Process all nodes at the current level
		for range size {
			cNode := q[0]
			q = q[1:] // POP Slice

			sum += cNode.Val

			// Add each child node to the queue
			if cNode.Left != nil {
				q = append(q, cNode.Left)
			}
			if cNode.Right != nil {
				q = append(q, cNode.Right)
			}
		}

		if sum > maxSum {
			maxSum = sum
			ans = level
		}
		level++

	}

	return ans

}
