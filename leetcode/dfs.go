package leetcode

func maxDepth(node *TreeNode) int {
	if node == nil {
		return 0
	}

	leftDepth := maxDepth(node.Left)
	rightDepth := maxDepth(node.Right)

	if leftDepth > rightDepth {
		return leftDepth + 1
	}
	return rightDepth + 1
}
