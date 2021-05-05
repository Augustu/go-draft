package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Augustu/go-draft/redis/values/types"
	"github.com/Augustu/go-draft/utils"
)

const (
	url string = "http://127.0.0.1:8888/values"
)

func postValues() {
	t1 := time.Now()

	ids := make(map[string]bool)

	// uids := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		id := utils.RandString(32)
		ids[id] = true

		v := &types.Value{
			ID:    id,
			Score: utils.RandomFloat64(),
		}
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println("marshal value failed", err)
			continue
		}

		body := bytes.NewBuffer(b)
		resp, err := http.Post(url, "", body)
		if err != nil || resp.StatusCode != http.StatusAccepted {
			fmt.Println(err)
		}

	}

	t2 := time.Now()
	fmt.Printf("ids: %d time: %s\n", len(ids), t2.Sub(t1).String())
}

func main() {
	for {
		postValues()
	}
}
