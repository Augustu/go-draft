package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
)

func chan1(c chan<- int) {
	defer func() {
		fmt.Println("write data done, exit chan1")
	}()
	for i := 0; i < 100005; i++ {
		c <- i
	}
	close(c)
}

func chan2(rc chan int, wc chan<- int, done chan<- bool) {

	defer func() {
		fmt.Println("exited chan2")
	}()

	for r := range rc {
		wc <- r
	}

	done <- true
	fmt.Println("rc chan2 closed")
}

func chan3(rc chan int, count *int64) {

	batchSize := 10
	var d []int

	defer func() {
		fmt.Println("exited chan3")
	}()

	for r := range rc {
		d = append(d, r)

		if len(d) >= batchSize {
			atomic.AddInt64(count, int64(len(d)))
			fmt.Println(d)
			d = []int{}
		}
	}

	atomic.AddInt64(count, int64(len(d)))

}

func main() {
	n := runtime.GOMAXPROCS(0)
	fmt.Println(n)

	c1 := make(chan int, 1000)

	go chan1(c1)

	c2 := make(chan int, 1000)

	c2chan := make(chan bool, n)
	for i := 0; i < n; i++ {
		c2chan <- true
	}

	for i := 0; i < n; i++ {
		<-c2chan
		go chan2(c1, c2, c2chan)
	}

	var count int64

	for i := 0; i < n; i++ {
		go chan3(c2, &count)
	}

	go func() {
		for i := 0; i < n; i++ {
			<-c2chan
		}
		close(c2)
		fmt.Println("all done")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-sc

	fmt.Println("\ncount", count)
}
