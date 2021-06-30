package main

import (
	"log"
	"net/http"
	"time"
)

const (
	url string = "http://127.0.0.1:8000/test"
)

func main() {
	for {
		w, err := http.Get(url)
		if err != nil {
			log.Println(err)
			break
		}
		w.Body.Close()
		time.Sleep(time.Millisecond)
	}
}
