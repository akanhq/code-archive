package hot_100

import (
	"sort"
	"testing"
)

func TestThreeSum(t *testing.T) {
	t.Logf("threeSum---result---> %v", threeSum([]int{-1, 0, 1, 2, -1, -4}))
}

func threeSum(nums []int) [][]int {
	result := [][]int{}

	if len(nums) < 3 {
		return result
	}

	sort.Ints(nums)

	return result
}
