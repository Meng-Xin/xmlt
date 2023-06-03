package utils

import "time"

// GetTimeSecond 获取秒
func GetTimeSecond() int {
	return time.Now().Second()
}

// GetTimeMilli 获取毫秒
func GetTimeMilli() int64 {
	return time.Now().UnixMilli()
}

// GetTimeNano 获取纳秒
func GetTimeNano() int {
	return time.Now().Nanosecond()
}
