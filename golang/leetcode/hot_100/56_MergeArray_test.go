package hot_100

import (
	"sort"
	"testing"
)

func TestMergeArray(t *testing.T) {

	t.Logf("MergeArray---result---> %v", merge([][]int{
		{1, 3},
		{2, 6},
		{8, 10},
		{15, 18},
	}))
}

func merge(intervals [][]int) [][]int {

	//按照起点重小到大排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := make([][]int, 0)
	curr := intervals[0]

	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] <= curr[1] {
			// 重叠,更新终点
			curr[1] = max(curr[1], intervals[i][1])
		} else {
			// 不重叠，加入结果,更新当前区间
			result = append(result, curr)
			curr = intervals[i]
		}
	}

	// 加入最后一个合并区间
	result = append(result, curr)

	return result
}
