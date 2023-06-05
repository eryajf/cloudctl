package public

import "time"

// cls 查询的结束时间
func GetNowTime() int64 {
	now := time.Now()
	return now.UnixNano() / int64(time.Millisecond)
}

// cls查询的开始时间，默认是往前查询3个小时
func GetOldTime() int64 {
	now := time.Now()
	threeHoursAgo := now.Add(-3 * time.Hour)
	return threeHoursAgo.UnixNano() / int64(time.Millisecond)
}

func Sum(nums []int) int {
	var total int
	for _, num := range nums {
		total += num
	}
	return total
}
