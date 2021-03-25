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

func Week() int {
	now := time.Now()

	future := now.AddDate(0, 0, 35)
	fmt.Println(future.Day())
	return 0
}
