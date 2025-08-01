package string_code

import (
	"fmt"
	"testing"
)

func longestCommonPrefix(strs []string) string {

	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]

	for i := 1; i < len(strs); i++ {
		j := 0
		for ; j < len(prefix) && j < len(strs[i]); j++ {
			if prefix[j] != strs[i][j] {
				break
			}
		}
		prefix = prefix[:j]
		if prefix == "" {
			return ""
		}
	}

	return prefix
}

func TestLongestCommonPrefix(t *testing.T) {
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
}
