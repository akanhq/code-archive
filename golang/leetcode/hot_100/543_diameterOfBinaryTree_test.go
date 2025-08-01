package hot_100

import (
	"testing"
)

func TestDiameterOfBinaryTree(t *testing.T) {}

func diameterOfBinaryTree(root *TreeNode) int {
	maxDiameter := 0
	var maxDepth func(node *TreeNode) int
	maxDepth = func(node *TreeNode) int {
		if node == nil {
			return 0
		}

		leftDepth := maxDepth(node.Left)
		rightDepth := maxDepth(node.Right)

		if leftDepth+rightDepth > maxDiameter {
			maxDiameter = leftDepth + rightDepth
		}
		return max(leftDepth, rightDepth) + 1
	}

	maxDepth(root)
	return maxDiameter
}
