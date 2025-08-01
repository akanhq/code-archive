package hot_100

import "testing"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func TestInorderTraversal(t *testing.T) {

}

func inorderTraversal1(root *TreeNode) []int {
	var result []int
	var stack []*TreeNode
	curr := root

	for curr != nil || len(stack) > 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}

		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, curr.Val)

		curr = curr.Right
	}
	return result
}
func inorderTraversal(root *TreeNode) []int {
	var result []int
	inorder(root, &result)
	return result
}

func inorder(node *TreeNode, result *[]int) {
	if node == nil {
		return
	}

	inorder(node.Left, result)
	*result = append(*result, node.Val)
	inorder(node.Right, result)
}
