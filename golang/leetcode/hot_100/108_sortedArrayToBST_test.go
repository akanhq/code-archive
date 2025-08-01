package hot_100

import "testing"

func TestSortedArrayToBST(t *testing.T) {

}

func sortedArrayToBST(nums []int) *TreeNode {

	return buildBST(nums, 0, len(nums)-1)
}

func buildBST(nums []int, left, right int) *TreeNode {

	if left > right {
		return nil
	}

	mid := left + (right-left)/2

	root := &TreeNode{Val: nums[mid]}

	root.Left = buildBST(nums, left, mid-1)
	root.Right = buildBST(nums, mid+1, right)

	return root
}
