package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

func resume(rc *ResumeClient) {
	ss, err := rc.CheckResume()
	if err != nil {
		fmt.Println("check resume failed", err)
		return
	}

	for _, s := range ss {
		fmt.Println("resume", s)
		work(rc, s.Index, s.Current)
		rc.Done(s.Index)
	}
}

func work(rc *ResumeClient, index int, current int) {
	err := rc.SetStatus(index, 10, "body-"+strconv.Itoa(index))
	if err != nil {
		fmt.Println("set status failed", index)
		return
	}

	go rc.Feed(index)

	for i := current; i < 10; i++ {
		rc.UpdateCurrent(index, i)

		fmt.Println(index, i)
		time.Sleep(time.Second)
	}
	rc.Done(index)
}

func main() {
	opts := redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	}
	rc := NewResumeClient(&opts, "resume:test")

	resume(rc)

	go work(rc, 0, 0)
	go work(rc, 1, 0)
	go work(rc, 2, 0)

	shutdownHandler := make(chan os.Signal, 1)
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

	signal.Notify(shutdownHandler, shutdownSignals...)
	<-shutdownHandler
	fmt.Println("shotdown now")
}
