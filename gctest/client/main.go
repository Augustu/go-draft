package main

import (
	"log"
	"net/http"
)

const (
	// url string = "http://127.0.0.1:8000/test"
	url string = "http://127.0.0.1:80/api/v1"

	// url string = "http://mission.rd-development.svc.cluster.local/api/v1"
)

func main() {
	for {
		// w, err := http.Get(url)
		_, err := http.Get(url)
		if err != nil {
			log.Println(err)
			break
		}
		// w.Body.Close()
		// time.Sleep(time.Millisecond)
	}
}
