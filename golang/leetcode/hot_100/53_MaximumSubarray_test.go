package hot_100

import (
	"testing"
)

func TestMaxSubArray(t *testing.T) {

	array := []int{1, 1, 1, 1, 1, 100}

	t.Logf("结果 %d", maxSubArray(array))
}

func maxSubArray1(nums []int) int {
	maxSum := nums[0]

	for i := 0; i < len(nums); i++ {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum > maxSum {
				maxSum = sum
			}
		}
	}

	return maxSum
}

func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	maxNum := nums[0]
	currSum := nums[0]

	for i := 1; i < len(nums); i++ {
		currSum = max(nums[i], currSum+nums[i])
		maxNum = max(currSum, maxNum)
	}

	return maxNum
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
