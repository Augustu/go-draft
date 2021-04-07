package timer

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func IntergralPoint() {
	now := time.Now()
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	fmt.Println(next.String())
}

func CronDemo() {
	c := cron.New()

	c.AddFunc("0 0-59/2 * * * *", func() {
		fmt.Println("Minutes test: ", time.Now().String())
	})

	c.Start()

	c.AddFunc("0 0-59/2 * * * *", func() {
		fmt.Println("Minutes test-2: ", time.Now().String())
	})

	num := 5
	spec := fmt.Sprintf("0 0-59/%d * * * *", num)
	c.AddFunc(spec, func() {
		fmt.Println("Minutes test-5: ", time.Now().String())
	})
	select {}
}

func WeekDemo() {
	c := cron.New()

	c.AddFunc("0 0 0 1 0-6/1", func() {
		fmt.Println("Week test: ", time.Now().String())
	})

	c.Start()

	select {}
}

func MinuteDemo() {
	c := cron.New()

	c.AddFunc("0 0-59/1 * * * *", func() {
		fmt.Println("Second test: ", time.Now().String())
	})

	c.Start()

	select {}
}

func SecondDemo() {
	c := cron.New()

	c.AddFunc("0-59/1 * * * * *", func() {
		fmt.Println("Second test: ", time.Now().String())
	})

	c.Start()

	select {}
}
