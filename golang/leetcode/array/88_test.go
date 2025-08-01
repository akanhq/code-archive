package main

import (
	"fmt"
	"testing"
)

//88. 合并两个有序数组

func merge(nums1 []int, m int, nums2 []int, n int) {

	if m == 0 {
		for i, i2 := range nums2 {
			nums1[i] = i2
		}
	}

	len1 := m - 1
	len2 := n - 1
	length := m + n - 1

	for len2 >= 0 && len1 >= 0 {
		if nums1[len1] > nums2[len2] {
			nums1[length] = nums1[len1]
			len1--
		} else {
			nums1[length] = nums2[len2]
			len2--
		}
		length--
	}

	fmt.Println("length", length)
}

func TestMerge(t *testing.T) {
	nums1 := []int{1, 2, 3, 0, 0, 0}
	nums2 := []int{2, 5, 6}
	merge(nums1, 3, nums2, 3)
	t.Log(nums1)
}
