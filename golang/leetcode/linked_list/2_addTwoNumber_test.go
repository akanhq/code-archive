package linked_list

import (
	"fmt"
	"strconv"
	"testing"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func arrToNumber(arr []int) string {

	numStr := ""

	for i := len(arr) - 1; i >= 0; i-- {
		numStr += strconv.Itoa(arr[i])
	}

	return numStr
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	dummy := &ListNode{} // 虚拟头节点
	current := dummy     // 当前节点
	carry := 0           // 进位

	// 遍历两个链表
	for l1 != nil || l2 != nil {
		// 获取当前节点的值
		var val1, val2 int
		if l1 != nil {
			val1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			val2 = l2.Val
			l2 = l2.Next
		}

		// 计算当前位的和
		sum := val1 + val2 + carry
		carry = sum / 10 // 更新进位
		current.Next = &ListNode{Val: sum % 10}
		current = current.Next
	}

	// 如果最后还有进位，添加一个新节点
	if carry > 0 {
		current.Next = &ListNode{Val: carry}
	}

	return dummy.Next
}

//func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
//
//	var l1i []int
//	var l2i []int
//
//	for l1 != nil {
//		l1i = append(l1i, l1.Val)
//		l1 = l1.Next
//	}
//
//	for l2 != nil {
//		l2i = append(l2i, l2.Val)
//		l2 = l2.Next
//	}
//
//	num1, _ := strconv.Atoi(arrToNumber(l1i))
//	num2, _ := strconv.Atoi(arrToNumber(l2i))
//
//	num := num1 + num2
//
//	numStr := strconv.Itoa(num)
//	runes := []rune(numStr)
//
//	val := &ListNode{}
//
//	for i := 0; i < len(runes); i++ {
//		val.Val, _ = strconv.Atoi(string(runes[i]))
//		ret := &ListNode{
//			Next: val,
//		}
//		val = ret
//	}
//	return val.Next
//}

func TestAddTwoNumbers(t *testing.T) {
	example1 := &ListNode{Val: 2, Next: &ListNode{Val: 4, Next: &ListNode{Val: 3, Next: nil}}}
	example2 := &ListNode{Val: 5, Next: &ListNode{Val: 6, Next: &ListNode{Val: 4, Next: nil}}}

	fmt.Println(addTwoNumbers(example1, example2).Val)

}
