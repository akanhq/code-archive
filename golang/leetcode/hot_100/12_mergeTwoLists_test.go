package hot_100

import (
	"testing"
)

func TestMergeTwoLists(t *testing.T) {
	tests := []struct {
		name     string
		list1    []int // 输入链表1
		list2    []int // 输入链表2
		expected []int // 预期输出
	}{
		{
			name:     "示例1: 正常合并",
			list1:    []int{1, 2, 4},
			list2:    []int{1, 3, 4},
			expected: []int{1, 1, 2, 3, 4, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建输入链表
			l1 := createList(tt.list1)
			l2 := createList(tt.list2)

			// 运行合并函数
			result := mergeTwoLists(l1, l2)
			println(result)

		})
	}
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	tempNode := &ListNode{}
	curr := tempNode

	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			curr.Next = list1
			list1 = list1.Next
		} else {
			curr.Next = list2
			list2 = list2.Next
		}
		curr = curr.Next
	}

	if list1 != nil {
		curr.Next = list1
	} else {
		curr.Next = list2
	}

	return tempNode.Next
}

// createList 辅助函数：从切片创建链表
func createList(nums []int) *ListNode {
	dummy := &ListNode{}
	current := dummy
	for _, val := range nums {
		current.Next = &ListNode{Val: val}
		current = current.Next
	}
	return dummy.Next
}
