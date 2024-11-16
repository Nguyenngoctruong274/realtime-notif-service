package xtime

import (
	"time"
)

const customRFC3339 = "2006-01-02T15:04:05"

func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

/*
*
Convert timestamp to string with custom RFC3339 format without GMT
ex: 2006-01-02T15:04:05Z07:00
*/
func ConvertUnixTimeRFC3339String(timeStamp int64) string {
	return time.Unix(timeStamp, 0).In(GMT07Location()).Format(customRFC3339)
}

/*
*
Convert timestamp to string with custom RFC3339 format without GMT
ex: 2006-01-02T15:04:05Z07:00
*/
func ParseStringToUnixTimestamp(timeStr string) int64 {
	t, err := time.Parse(customRFC3339, timeStr)
	if err != nil {
		return 0
	}
	return t.Unix()
}

/*
*
Convert timestamp to string with custom RFC3339 format without GMT
ex: 2006-01-02T15:04:05Z07:00
*/
func ParseStringToUnixTimestampLocation(timeStr string) int64 {
	t, err := time.ParseInLocation(customRFC3339, timeStr, GMT07Location())
	if err != nil {
		return 0
	}
	return t.Unix()
}
