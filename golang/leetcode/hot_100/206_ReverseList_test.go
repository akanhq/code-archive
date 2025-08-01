package hot_100

import "testing"

func TestReverseList(t *testing.T) {

	t1 := &ListNode{
		Val: 1,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 3,
				Next: &ListNode{
					Val: 4,
					Next: &ListNode{
						Val:  5,
						Next: nil,
					},
				},
			},
		},
	}

	t.Logf("reverseList---result---> %v", reverseList(t1))
}

func reverseList(head *ListNode) *ListNode {
	result := &ListNode{}
	curr := head

	if curr != nil {
		for curr != nil { // 遍历链表，直到 curr 为 nil
			next := curr.Next  // 保存下一个节点，避免丢失
			curr.Next = result // 反转：将 curr 的 Next 指向 prev
			result = curr      // 更新 prev 为当前节点（新的反转头节点）
			curr = next        // 移动 curr 到下一个节点
		}
	}
	return result
}
