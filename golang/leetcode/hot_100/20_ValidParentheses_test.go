package hot_100

import (
	"fmt"
	"testing"
)

func TestValidParentheses(t *testing.T) {
	fmt.Println(ValidParentheses("([])"))
}

func ValidParentheses(s string) bool {
	stack := []rune{}

	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char)
			continue
		}
		if len(stack) == 0 {
			return false
		}

		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if (char == ')' && top != '(') || (char == '}' && top != '{') || (char == ']' && top != '[') {
			return false
		}
	}
	return len(stack) == 0
}
