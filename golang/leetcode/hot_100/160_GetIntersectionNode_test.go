package hot_100

import "testing"

type ListNode struct {
	Val  int
	Next *ListNode
}

func TestGetIntersectionNode(t *testing.T) {}
func getIntersectionNode(headA, headB *ListNode) *ListNode {
	pa, pb := headA, headB
	for pa != pb {
		if pa != nil {
			pa = pa.Next
		} else {
			pa = headB
		}

		if pb != nil {
			pb = pb.Next
		} else {
			pb = headA
		}
	}
	return pa
}
