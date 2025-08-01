package hot_100

import "testing"

func TestProductExceptSelf(t *testing.T) {
	t.Logf("rotateArray---result---> %v", productExceptSelf([]int{1, 2, 3, 4}))
}

func productExceptSelf(nums []int) []int {
	n := len(nums)
	answer := make([]int, n)

	// 计算左侧乘积
	leftProd := 1
	for i := 0; i < n; i++ {
		answer[i] = leftProd
		leftProd *= nums[i]
	}

	//计算右侧乘积
	rightProd := 1
	for i := n - 1; i >= 0; i-- {
		answer[i] *= rightProd
		rightProd *= nums[i]
	}

	return answer
}
