package util

import "time"

const (
	TZDateTimeFormat  = "2006-01-02T15:04:05Z"
	UTCDateTimeFormat = "2006-01-02T15:04:05Z07:00"
	DateTimeFormat    = "2006-01-02 15:04:05"
	DateFormat        = "2006-01-02"
	HourFormat        = "15:04"
)

// HourTime 时间格式转换
func HourTime(time time.Time) string {
	return time.Format(HourFormat)
}

// DateTime 时间格式转换
func DateTime(time time.Time) string {
	return time.Format(DateTimeFormat)
}

// UTCDateTime 时间格式转换(UTC)
func UTCDateTime(time time.Time) string {
	return time.Format(UTCDateTimeFormat)
}

// Date 日期格式转换
func Date(time time.Time) string {
	return time.Format(DateFormat)
}

// TimeFromString 字符串转换为时间
func TimeFromString(datetime string) (time.Time, error) {
	return time.Parse(DateFormat, datetime)
}

// TimeFromTZDateString 时区字符串，转时间
func TimeFromTZDateString(datetime string) (time.Time, error) {
	return time.Parse(TZDateTimeFormat, datetime)
}

// TimeFromUTCString UTC字符串转换为时间
func TimeFromUTCString(datetime string) (time.Time, error) {
	return time.Parse(UTCDateTimeFormat, datetime)
}

// TimeStampAfterSeconds 返回多少秒后的时间戳
func TimeStampAfterSeconds(seconds int64) int64 {
	return time.Now().Unix() + seconds
}

// TimeStampAfterMillSeconds 返回 多少秒后的 毫秒时间戳
func TimeStampAfterMillSeconds(seconds int64) int64 {
	return time.Now().UnixMilli() + (seconds * 1000)
}

// TimeFromMillSeconds 毫秒时间戳 转换成时间
func TimeFromMillSeconds(ms int64) time.Time {
	return time.UnixMilli(ms)
}

// GetNow 获取当前时间
func GetNow() time.Time {
	return time.Now()
}

// GetYear 获取年份
func GetYear(now *time.Time) int {
	return now.Year()
}

// GetMonth 获取月份
func GetMonth(now *time.Time) int {
	return int(now.Month())
}

// GetDay 获取日期
func GetDay(now *time.Time) int {
	return now.Day()
}

func GetHour(now *time.Time) int {
	return now.Hour()
}

func GetMin(now *time.Time) int {
	return now.Minute()
}
