package hot_100

import (
	"testing"
)

func TestRotateArray(t *testing.T) {
	req := []int{1, 2}
	rotate(req, 3)
	t.Logf("rotateArray---result---> %v", req)
}

func rotate(nums []int, k int) {
	if k == 0 {
		return
	}

	k = k % len(nums)
	lastIndex := len(nums) - k

	result := nums[lastIndex:]
	for i := 0; i < len(nums)-k; i++ {
		result = append(result, nums[i])
	}

	copy(nums, result)
}
