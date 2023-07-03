package utils

import (
	"math/rand"
	"time"
)

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

// GetRandTime 获取指定范围时间
// base 基础时间值 append 追加时间 max 随机偏移值
func GetRandTime(baseTime, appendUnit time.Duration, max int) time.Duration {
	return baseTime + (appendUnit * time.Duration(rand.Intn(max)))
}
