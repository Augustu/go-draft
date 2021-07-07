package main

import "time"

func main() {
	s := []string{"1", "2", "3"}
	for {
		s = append(s, s...)
		time.Sleep(100 * time.Millisecond)
		if len(s) > 10e16 {
			break
		}
	}
}
