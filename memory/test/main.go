package main

import "time"

func test() {
	s := []string{"1", "2", "3"}
	for {
		s = append(s, s...)
		time.Sleep(100 * time.Millisecond)
		if len(s) > 10e16 {
			break
		}
	}
}

func testBigChan() {
	ch := make(chan string, 1024000)

	for i := 0; i < 1024000; i++ {
		ch <- "test"
	}
}

func main() {
	testBigChan()
}
