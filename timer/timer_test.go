package timer

import (
	"fmt"
	"testing"
	"time"
)

// func TestIntergralPoint(t *testing.T) {
// 	IntergralPoint()
// }

// func TestCronDemo(t *testing.T) {
// 	CronDemo()
// }

// func TestWeekDemo(t *testing.T) {
// 	WeekDemo()
// }

// func TestMinuteDemo(t *testing.T) {
// 	MinuteDemo()
// }

// func TestSecondDemo(t *testing.T) {
// 	SecondDemo()
// }

func TestTimer(t *testing.T) {
	tr := time.NewTimer(3 * time.Second)
	<-tr.C
	fmt.Println("timeout")

}
