package hot_100

import "testing"

func TestMoveZeroes(t *testing.T) {
	req := []int{0, 1, 0, 3, 12}
	moveZeroes(req)
	t.Logf("moveZeroes---result---> %v", req)
}

func moveZeroes(nums []int) {
	nonZeroPos := 0

	for i := 0; i < len(nums); i++ {
		if nums[i] != 0 {
			nums[nonZeroPos] = nums[i]
			nonZeroPos++
		}
	}

	for i := nonZeroPos; i < len(nums); i++ {
		nums[i] = 0
	}
}
