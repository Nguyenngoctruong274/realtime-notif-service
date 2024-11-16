package timeutils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"yes4all/ads-noti-api/pkg/utils/constants"
)

const (
	customRFC3339    = "2006-01-02T15:04:05"
	ISO8601Layout    = "2006-01-02T15:04:05Z0700"
	ISOLayout        = "2006-01-02"
	hourInDay        = 24
	CustomDateTimeVN = "15:04:05 02-01-2006"
	DateTime         = "2006-01-02 15:04:05"
)

var (
	locationGMT07 *time.Location
	locationLA    *time.Location
	once          sync.Once
	initialized   bool
)

type NowFn func() time.Time
type NowTimestampFn func() int64

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
		locationLA, err = time.LoadLocation("America/Los_Angeles")
		if err != nil {
			panic(err)
		}
		initialized = true
	})
}

func gmt07Location() *time.Location {
	if !initialized {
		fmt.Println("cannot use GMT+07 timezone, have you forgotten to call InitTimezones()?")
		return time.UTC
	}
	return locationGMT07
}

func americaLALocaltion(ctx context.Context) *time.Location {
	if !initialized {
		fmt.Println("cant initialize value for timezone LA")
		return time.UTC
	}
	location, _ := time.LoadLocation("America/ Los_Angeles")
	return location
}

var NowInGMT07 = func() time.Time {
	return ConvertTimeToGMT07(time.Now())
}

func TimeInGMT07String(t time.Time) string {
	result := t.In(gmt07Location()).Format(customRFC3339)
	return result
}

func ConvertTimestampToTime(timestamp int64) time.Time {
	return ConvertTimeToGMT07(time.Unix(timestamp, 0))
}

func ConvertTimeToGMT07(t time.Time) time.Time {
	return t.In(gmt07Location())
}

func ConvertTimeToAMZ(ctx context.Context, t time.Time) time.Time {
	return t.In(americaLALocaltion(ctx))
}

func ConvertTimeStampToString(timeStamp int64, format string) string {
	t := ConvertTimestampToTime(timeStamp)
	return t.Format(format)
}

func ConvertTimeToString(t time.Time, format string) string {
	t = ConvertTimeToGMT07(t)
	return t.Format(format)
}

func IsWeekend(datetime time.Time, includeSaturday bool) bool {
	dayOfWeek := datetime.Weekday()
	return dayOfWeek == time.Sunday || (dayOfWeek == time.Saturday && includeSaturday)
}

func BeginningOfDay(datetime time.Time) time.Time {
	y, m, d := datetime.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, gmt07Location())
}

func BeginningOfMonth(datetime time.Time) time.Time {
	y, m, _ := datetime.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, gmt07Location())
}

func BeginningOfMonthLAZone(ctx context.Context, datetime time.Time) time.Time {
	y, m, _ := datetime.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, americaLALocaltion(ctx))
}

func CompareByDate(timeA time.Time, timeB time.Time) bool {
	return timeA.Truncate(hourInDay * time.Hour).Equal(timeB.Truncate(hourInDay * time.Hour))
}

func EndOfDayGMT07(datetime time.Time) time.Time {
	y, m, d := datetime.Date()
	return time.Date(y, m, d, 23, 59, 59, 0, gmt07Location())
}

func EndOfDay(datetime time.Time) time.Time {
	y, m, d := datetime.Date()
	return time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
}

func ConvertTimeNilToTimeStamp(datetime *time.Time) int64 {
	if datetime == nil {
		return 0
	}
	result := ConvertTimeToTimeStamp(*datetime)
	return result
}

// Trường hợp datetime truyền vào là UTC ( thường là các datetime từ database lên ).
// Hệ thống sẽ convert qua GMT7 --> parse sang timestamp
func ConvertTimeToTimeStamp(datetime time.Time) int64 {
	if datetime.IsZero() {
		return 0
	}
	name, _ := datetime.Zone()
	if name == "UTC" {
		datetime = ConvertWithoutTimeZoneToGMT7Time(datetime)
	}
	result := datetime.Unix()
	return result
}

func Since(source time.Time) time.Duration {
	return NowInGMT07().Sub(source)
}

// JSONDateTime first create a type alias
type JSONDateTime time.Time

var layoutsWithTimeZone = []string{time.RFC3339, time.RFC3339Nano, ISO8601Layout}
var layoutsWithoutTimeZone = []string{
	customRFC3339,
	"2006-01-02T15:04:05.999999",
	"2006-01-02 15:04:05",
	constants.TimeLayoutddMMyyyytimePattern,
	ISOLayout,
}

func ConvertStringToTime(src string) time.Time {
	if len(src) < 1 {
		return time.Time{}
	}
	for _, layout := range layoutsWithTimeZone {
		t, err := time.Parse(layout, src)
		if err == nil {
			return t
		}
	}
	for _, layout := range layoutsWithoutTimeZone {
		t, err := time.ParseInLocation(layout, src, gmt07Location())
		if err == nil {
			return t
		}
	}
	return time.Time{}
}

func ConvertStringToTimeAMZ(ctx context.Context, src string, layout string) time.Time {
	if len(src) < 1 {
		return time.Time{}
	}

	t, err := time.ParseInLocation(layout, src, americaLALocaltion(ctx))
	if err == nil {
		return t
	}
	return time.Time{}
}

func ConvertStringToTimeFormat(src string, format string) time.Time {
	if len(src) < 1 {
		return time.Time{}
	}
	t, err := time.Parse(format, src)
	if err == nil {
		return t
	}
	return time.Time{}
}

func ValidateTimeString(dateString string, format string) (err error) {
	_, err = time.Parse(format, dateString)
	if err != nil {
		return
	}

	return nil
}

func (j JSONDateTime) MarshalJSON() ([]byte, error) {
	/*stamp := fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02T15:04:05.999999"))
	return []byte(stamp), nil*/
	return []byte(j.String()), nil
}

func (j *JSONDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t := ConvertStringToTime(s)
	if t.IsZero() {
		return errors.New("not support datetime layout")
	}
	*j = JSONDateTime(t)
	return nil
}

// String returns the time in the custom format
func (j *JSONDateTime) String() string {
	t := time.Time(*j)
	return fmt.Sprintf("%q", t.Format(time.RFC3339))
}

func ConvertJSONTimeToTime(j JSONDateTime) time.Time {
	return (time.Time)(j)
}

func ConvertWithoutTimeZoneToGMT7Time(t time.Time) time.Time {
	dateStr := t.Format("2006-01-02 15:04:05")
	return ConvertStringToTime(dateStr)
}

func ConvertWithoutTimeZoneToGMT7Date(t time.Time) time.Time {
	dateStr := t.Format(ISOLayout)
	return ConvertStringToTime(dateStr)
}

func ConvertTimeStringToNewLayout(layout string, t time.Time) string {
	return t.Format(layout)
}

func ConvertWithoutTimeZoneToTimestamp(t time.Time) int64 {
	t = ConvertWithoutTimeZoneToGMT7Time(t)
	return ConvertTimeToTimeStamp(t)
}

func SetDateTime(y int, m time.Month, d int, hour int, min int, sec int) time.Time {
	return time.Date(y, m, d, hour, min, sec, 0, gmt07Location())
}

func GetStartDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, gmt07Location())
}

func GetEndDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 59, gmt07Location())
}

func GetStartMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, gmt07Location())
}

func GetEndMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 0, gmt07Location())
}

func GetEndMonthLA(ctx context.Context, date time.Time) time.Time {
	return time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 0, americaLALocaltion(ctx))
}

func GetStartYear(date time.Time) time.Time {
	return time.Date(date.Year(), 1, 1, 0, 0, 0, 0, gmt07Location())
}

func TruncateSecondTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, gmt07Location())
}

func ConvertTimeToTimePointer(t time.Time) *time.Time {
	return &t
}

// Nghiệp vụ sẽ thay đổi theo yêu cầu chu kỳ loyalty của biz.
// Hiện tại đang quy định mỗi chu kỳ là 1 tháng.
func GenLoyaltyPeriodByDate(t time.Time) (startAt, endAt time.Time) {
	startAt = BeginningOfMonth(t)
	endAt = EndOfDayGMT07(startAt.AddDate(0, 1, -1))
	return
}

func GetMonth(timestamp int64) int {
	t := ConvertTimestampToTime(timestamp)
	return int(t.Month())
}

func NowWithTruncateOneDay() time.Time {
	t := NowInGMT07()
	someMin := (t.Minute() / constants.PartOfTime) * constants.PartOfTime
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), someMin, 0, 0, gmt07Location())
}

func IsNextDay(timestamp int64) bool {
	// return timestamp-NowWithTruncateOneDay().Unix() < 0 // FOR TEST
	return timestamp-BeginningOfDay(NowInGMT07()).Unix() < 0
}

func NewDayWithTruncateOneDay(t time.Time) time.Time {
	someMin := (t.Minute() / constants.PartOfTime) * constants.PartOfTime
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), someMin, 0, 0, gmt07Location())
}

func GenDate(startTime, endTime time.Time, layout string) (res []string) {
	start := startTime
	for start.Format("2006-01-02") <= endTime.Format("2006-01-02") {
		res = append(res, start.Format(layout))
		start = start.AddDate(0, 0, 1)
	}
	return
}
