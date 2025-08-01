package hot_100

import (
	"fmt"
	"testing"
)

func TestHasCycle(t *testing.T) {
	fmt.Println(hasCycle(createList([]int{1, 2})))
}

func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow, fast := head, head.Next
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
		if fast == slow {
			return true
		}
	}
	return false
}
