package main

import "fmt"

func chan1(c chan<- int, done chan bool) {
	defer func() {
		fmt.Println("write data done, exit chan1")
	}()
	for i := 0; i < 105; i++ {
		c <- i
	}
	close(done)
}

func main() {

}
