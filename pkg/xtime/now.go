package xtime

import "time"

type NowInMillisecond func() int64
type NowTimeFn func() time.Time

// Now returns unix timestamp in milliseconds
func NowMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func NowSec() int64 {
	return time.Now().Unix()
}

func NowTime() time.Time {
	return time.Now()
}
