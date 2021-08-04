package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	TimeLayout string = "2006-01-02 15:04:05"
)

func RandomInt(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(end - start)
	random = start + random
	return random
}

func RandomFloat64() float64 {
	return rand.Float64()
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func RandomTime(startYear, endYear, startMonth, endMonth, startDay, endDay, startHour, endHour, startMinute, endMinute, startSecond, endSecond int) time.Time {
	y := RandomInt(startYear, endYear)
	m := RandomInt(startMonth, endMonth)
	d := RandomInt(startDay, endDay)
	h := RandomInt(startHour, endHour)
	min := RandomInt(startMinute, endMinute)
	s := RandomInt(startSecond, endSecond)

	ts := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", y, m, d, h, min, s)
	t, err := time.Parse(TimeLayout, ts)
	if err != nil {
		return time.Time{}
	}

	return t
}
