package hot_100

import (
	"fmt"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	fmt.Println(isPalindrome(createList([]int{1, 2, 2, 1})))
}
func isPalindrome(head *ListNode) bool {

	var list []int
	for head != nil {
		list = append(list, head.Val)
		head = head.Next
	}

	i, j := 0, len(list)-1

	for i < j {
		if list[i] != list[j] {
			return false
		}
		j--
		i++
	}
	return true

}
