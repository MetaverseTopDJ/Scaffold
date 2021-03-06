package util

import "time"

const (
	TZDateTimeFormat        = "2006-01-02T15:04:05Z"
	UTCDateTimeFormat       = "2006-01-02T15:04:05Z07:00"
	UTCDateTimeMinuteFormat = "2006-01-02T15:04Z07:00"
	DateTimeFormat          = "2006-01-02 15:04:05"
	DateFormat              = "2006-01-02"
	DateMinuteFormat        = "2006-01-02 15:04"
	HourFormat              = "15:04"
	DateTimeStringFormat    = "20060102150405"
)

// StartEndFromString 时间字符串 转换为时间
func StartEndFromString(s, e, format string) (start, end time.Time, err error) {
	start, err = time.Parse(format, s)
	if err != nil {
		return
	}
	end, err = time.Parse(format, e)
	if err != nil {
		return
	}
	return
}

// TimeToString 时间转换为字符串
func TimeToString(t time.Time, f string) string {
	return t.Format(f)
}

// HourTime 时间格式转换
func HourTime(time time.Time) string {
	return time.Format(HourFormat)
}

// DateTime 时间格式转换
func DateTime(time time.Time) string {
	return time.Format(DateTimeFormat)
}

// DateTimeString 时间格式转换
func DateTimeString(time time.Time) string {
	return time.Format(DateTimeStringFormat)
}

// DateMinuteTime 时间格式转换
func DateMinuteTime(time time.Time) string {
	return time.Format(DateMinuteFormat)
}

// UTCDateTime 时间格式转换(UTC)
func UTCDateTime(time time.Time) string {
	return time.Format(UTCDateTimeFormat)
}

// UTCDateMinuteTime 时间格式转换(UTC)
func UTCDateMinuteTime(time time.Time) string {
	return time.Format(UTCDateTimeMinuteFormat)
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

// TimeMinuteFromUTCString UTC字符串转换为时间(精确到分)
func TimeMinuteFromUTCString(datetime string) (time.Time, error) {
	return time.Parse(UTCDateTimeMinuteFormat, datetime)
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

// StringFromTimestamp 时间戳 转换为 字符串
func DateTimeFromTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(DateTimeFormat)
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
