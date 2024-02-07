package util

import (
	"time"
)

func DayBegin() time.Time {
	end_ := time.Now()
	_, offset := end_.Zone()
	s := end_.Unix()
	s = s - (s+int64(offset))%(24*60*60)
	return time.Unix(s, 0)
}

// NowMs 当前时间戳（毫秒）
func NowMs() int64 {
	return time.Now().UnixNano() / 1e6
}
