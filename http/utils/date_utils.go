package utils

import "time"

const (
	FORMATTER_DEFAULT = "2006-01-02 15:04:05"
)

func FormatNow() (timeStr string) {
	return time.Now().Format(FORMATTER_DEFAULT)
}

func FormatNowBy(formatter string) (timeStr string) {
	return time.Now().Format(formatter)
}

func FormatTime(time time.Time) (timeStr string) {
	return time.Format(FORMATTER_DEFAULT)
}

func FormatTimeBy(time time.Time, formatter string) (timeStr string) {
	return time.Format(formatter)
}

func ParseTime(timeStr string) (timeObj time.Time, err error) {
	return time.Parse(FORMATTER_DEFAULT, timeStr)
}

func ParseTimeBy(timeStr string, formatter string) (timeObj time.Time, err error) {
	return time.Parse(formatter, timeStr)
}