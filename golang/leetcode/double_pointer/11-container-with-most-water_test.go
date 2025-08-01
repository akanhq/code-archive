package double_pointer

import (
	"fmt"
	"testing"
)

func maxArea(height []int) int {

	var max int
	i, j := 0, len(height)-1
	lenght := len(height) - 1

	for i < j {
		if height[i] > height[j] {
			result := height[j] * lenght
			if result > max {
				max = result
			}
			j--
		} else {
			result := height[i] * lenght
			if result > max {
				max = result
			}
			i++
		}
		lenght--
	}

	return max
}

func TestContainerWithMostWater(t *testing.T) {
	fmt.Println(maxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7}))
}
