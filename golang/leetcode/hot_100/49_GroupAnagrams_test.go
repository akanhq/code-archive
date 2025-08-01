package hot_100

import (
	"sort"
	"testing"
)

func TestGroupAnagrams(t *testing.T) {
	t.Logf("result = %v", groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat", "bat"}))
}

func groupAnagrams(strs []string) [][]string {

	m := make(map[string][]string)

	for _, str := range strs {
		chars := []byte(str)
		sort.Slice(chars, func(i, j int) bool {
			return chars[i] < chars[j]
		})

		sortedStr := string(chars)
		m[sortedStr] = append(m[sortedStr], str)
	}

	result := make([][]string, 0, len(m))
	for _, group := range m {
		result = append(result, group)
	}

	return result
}
