package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Augustu/go-draft/cache/cmd/server/model"
	"github.com/Augustu/go-draft/utils"
)

const (
	addr string = "http://127.0.0.1:8000/create"
)

func post(s model.Stats) {
	body, err := json.Marshal(s)
	if err != nil {
		return
	}

	_, err = http.Post(addr, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
}

func postBatch(statsChan <-chan model.Stats, n int) {
	threads := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		threads <- struct{}{}
	}

	for s := range statsChan {
		<-threads
		go func(s model.Stats) {
			post(s)
			threads <- struct{}{}
		}(s)
	}
}

func genStats(statsChan chan<- model.Stats, n int) {
	for i := 0; i < n; i++ {
		statsChan <- model.Stats{
			A:          utils.RandomInt(0, 10),
			B:          utils.RandomInt(0, 1000),
			C:          utils.RandomInt(0, 10000),
			D:          utils.RandomInt(0, 1000000),
			OccurredAt: utils.RandomTime(2021, 2022, 5, 7, 1, 30, 0, 23, 0, 59, 0, 59),
		}
	}
	close(statsChan)
}

func main() {
	n := 10000000
	threads := 8

	statsChan := make(chan model.Stats, 100000)

	go genStats(statsChan, n)
	postBatch(statsChan, threads)
}
