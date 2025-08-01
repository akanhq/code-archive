package main

import (
	"fmt"
	"testing"
)

func singleNumber(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}

	var numMap = make(map[int][]int)
	for _, num := range nums {
		numMap[num] = append(numMap[num], num)
	}

	for key, vals := range numMap {
		if len(vals)%2 != 0 {
			return key
		}
	}
	return 0

}
func TestSingleNumber(t *testing.T) {
	prices := []int{2, 4, 1}
	fmt.Printf("只出现一次的数字：%d", singleNumber(prices))
}
