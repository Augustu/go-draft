package time

import (
	"fmt"
	"time"
)

func String() string {
	now := time.Now()
	// return fmt.Sprintf("%04d-%02d-%02d-%02d-%02d-%02d",
	// 	now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
	// )
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
	)
}

// func WeekTime(ts string) (time.Time, error) {
// 	if len(ts) != 6 {
// 		return time.Time{}, fmt.Errorf("invalid week time string: %s", ts)
// 	}

// 	// t, err := time.Parse("")
// 	if err != nil {
// 		return time.Time{}, fmt.Errorf("invalid week time year string: %s", ts[0:4])
// 	}

// 	w, err := strconv.Atoi(ts[3:5])
// 	if err != nil {
// 		return time.Time{}, fmt.Errorf("invalid week time weak string: %s", ts[3:5])
// 	}
// }

func MinuteString(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func HourString(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour())
}

func DayString(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

func WeekString(t time.Time) string {
	y, w := t.ISOWeek()
	return fmt.Sprintf("%04d%02d", y, w)
}

func MonthString(t time.Time) string {
	return fmt.Sprintf("%04d%02d", t.Year(), t.Month())
}

func YearIntString(y int) string {
	return fmt.Sprintf("%04d", y)
}

func YearString(t time.Time) string {
	return fmt.Sprintf("%04d", t.Year())
}

func LastHourString(t time.Time) string {
	duration, _ := time.ParseDuration("1h")
	t = t.Add(duration)
	return fmt.Sprintf("%04d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour())
}

func LastDayString(t time.Time) string {
	t = t.Add(-24 * time.Hour)
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

func LastWeekString(t time.Time) string {
	t = t.Add(-7 * 24 * time.Hour)
	y, w := t.ISOWeek()
	return fmt.Sprintf("%04d%02d", y, w)
}

func LastMonthString(t time.Time) string {
	t = t.Add(-30 * 24 * time.Hour)
	return fmt.Sprintf("%04d%02d", t.Year(), t.Month())
}

func Hours(start time.Time, end time.Time) []string {
	var hours []string
	for s := start; !s.After(end); s.Add(time.Hour) {
		hours = append(hours, HourString(s))
	}
	return hours
}

func Days(start time.Time, end time.Time) []string {
	var days []string
	for s := start; !s.After(end); s.Add(24 * time.Hour) {
		days = append(days, DayString(s))
	}
	return days
}

func Weeks(start time.Time, end time.Time) []string {
	var weeks []string
	for s := start; !s.After(end); s.Add(7 * 24 * time.Hour) {
		weeks = append(weeks, WeekString(s))
	}
	return weeks
}

func Months(start time.Time, end time.Time) []string {
	var months []string
	// use 30 day for month is enough
	for s := start; !s.After(end); s.Add(30 * 24 * time.Hour) {
		months = append(months, WeekString(s))
	}
	return months
}

func Years(start time.Time, end time.Time) []string {
	var years []string
	for s := start.Year(); s <= end.Year(); s++ {
		years = append(years, YearIntString(s))
	}
	return years
}

func Week() int {
	now := time.Now()

	future := now.AddDate(0, 0, 35)
	fmt.Println(future)
	fmt.Println(future.Round(time.Hour))
	fmt.Println(future.Round(24 * time.Hour))
	fmt.Println(future.Round(7 * 24 * time.Hour))
	fmt.Println(future.Round(30 * 24 * time.Hour))
	return 0
}
