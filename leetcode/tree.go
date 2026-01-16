package leetcode

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type treeTest struct {
	tree   *TreeNode
	expect int
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
	fmt.Println("\n", arr)
	printTreeStructure("ROOT", root)
	return root
}

func printTreeStructure(label string, root *TreeNode) {

	if root == nil {
		return
	}

	leftVal := "nil"
	if root.Left != nil {
		leftVal = fmt.Sprintf("%d", root.Left.Val)
	}
	rightVal := "nil"
	if root.Right != nil {
		rightVal = fmt.Sprintf("%d", root.Right.Val)
	}

	fmt.Printf("Left: %s ------ Parent-%s: %v ----- Right:%s \n", leftVal, label, root.Val, rightVal)

	printTreeStructure("Left", root.Left)
	printTreeStructure("Right", root.Right)

}
