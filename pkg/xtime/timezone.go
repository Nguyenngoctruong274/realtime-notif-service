package xtime

import (
	"sync"
	"time"
	"yes4all/ads-noti-api/pkg/logger"
)

var (
	locationGMT07 *time.Location
	once          sync.Once
	initialized   bool
)

//nolint:gochecknoinits
func init() {
	initTimezones()
}

func initTimezones() {
	once.Do(func() {
		var err error
		// Load required location
		locationGMT07, err = time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			panic(err)
		}
		initialized = true
	})
}

// InitTimezones initializes GMT07 timezone data
// Deprecated: the clock pkg will initialize timezone data automatically, no need to init explicitly
func InitTimezones() {
	initTimezones()
}

func logWarning() {
	logger.NewLogger().Warn("cannot use GMT+07 timezone, have you forgotten to call InitTimezones()?")
}

func GMT07Location() *time.Location {
	if !initialized {
		logWarning()
		return time.UTC
	}
	return locationGMT07
}

func NowInBeautifiedGMT07String() string {
	return NowInGMT07String("02.01.2006 15:04")
}

func NowInGMT07String(format string) string {
	return TimeInGMT07String(time.Now(), format)
}

func TimeInBeautifiedGMT07String(t time.Time) string {
	return TimeInGMT07String(t, "02.01.2006 15:04")
}

func TimeInGMT07String(t time.Time, format string) string {
	result := t.In(GMT07Location()).Format(format)
	return result
}

func TimestampToGMT07Time(timestamp int64) *time.Time {
	timeGMT07 := ConvertTimeToGMT07(time.Unix(timestamp, 0))
	return &timeGMT07
}

func ConvertTimeToGMT07(t time.Time) time.Time {
	return t.In(GMT07Location())
}
