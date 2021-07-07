package utils

import "time"

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation(TIME_FORMAT, s, time.Local)
}
