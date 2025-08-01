package hot_100

import (
	"testing"
)

func TestLongestConsecutive(t *testing.T) {
	t.Logf("longestConsecutive---result---> %v", longestConsecutive([]int{100, 4, 200, 1, 3, 2}))
}

func longestConsecutive(nums []int) int {

	numSet := make(map[int]struct{})
	for _, num := range nums {
		numSet[num] = struct{}{}
	}
	maxLen := 0

	//找到起点
	for num, _ := range numSet {
		if _, exists := numSet[num-1]; exists {
			continue
		}

		currNum := num
		currLen := 1
		//往后连续查找
		_, exists := numSet[currNum+1]
		for exists {
			currNum++
			currLen++
			_, exists = numSet[currNum+1]
		}

		// 更新最大长度
		if currLen > maxLen {
			maxLen = currLen
		}
	}

	return maxLen
}
