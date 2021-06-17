package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runtine(done <-chan bool, id int) {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("tick", id)

		case <-done:
			fmt.Println("done", id)
			return
		}
	}
}

func main() {
	done := make(chan bool)

	for i := 0; i < 3; i++ {
		go runtine(done, i)
	}

	time.Sleep(5 * time.Second)
	close(done)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-sc
}
