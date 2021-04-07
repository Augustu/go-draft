package time

import (
	"fmt"
	"testing"
	"time"
)

// func TestString(t *testing.T) {
// 	fmt.Println(String())
// }

// func TestWeek(t *testing.T) {
// 	Week()
// }

func TestTime(t *testing.T) {
	n := time.Now()
	fmt.Println(n.Format("2006-01-02 15:04:05"))

	tt := "2021-04-07 14:40:09"

	pt, _ := time.ParseInLocation("2006-01-02 15:04:05", tt, time.Local)
	fmt.Println(n)
	fmt.Println(pt)
}
