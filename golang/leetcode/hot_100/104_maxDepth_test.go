package hot_100

import "testing"

func TestMaxDepth(t *testing.T) {

}

func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	leftDepth := maxDepth(root.Left)
	rightDepth := maxDepth(root.Right)
	return max(leftDepth, rightDepth) + 1
}
