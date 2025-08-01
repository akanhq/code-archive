package hot_100

import (
	"fmt"
	"testing"
)

func TestSubarraySum(t *testing.T) {
	fmt.Printf("subarraySum %d", subarraySum([]int{1, 1, 1}, 2))
}
func subarraySum(nums []int, k int) int {
	count := 0

	n := len(nums)

	for i := 0; i < n; i++ {

		sum := 0

		for j := i; j < n; j++ {
			sum += nums[j]
			if sum == k {
				count++
			}
		}
	}
	return count

}
