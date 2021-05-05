package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Augustu/go-draft/redis/values/types"
	"github.com/go-redis/redis/v8"
)

var (
	key string = "zset"
	c   *redis.Client

	valuesChan = make(chan types.Value, 102400)
)

func valuesHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var v types.Value
	err = json.Unmarshal(body, &v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	valuesChan <- v
}

func cacheValues(ctx context.Context) {
	var count int64

	for {
		select {
		case <-ctx.Done():
			fmt.Println("exit cache values")
			return

		case v := <-valuesChan:
			r, e := c.ZAdd(ctx, key, &redis.Z{
				Score:  v.Score,
				Member: v.ID,
			}).Result()
			if e != nil {
				fmt.Println(r, e)
			}
			count++
			fmt.Println(count)
		}
	}
}

func snapshot(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("exit snapshot")

		case <-ticker.C:
			r, err := c.Rename(ctx, key, fmt.Sprintf("%s-%s", key, time.Now().String())).Result()
			if err != nil {
				fmt.Println(r, err)
			}
		}
	}
}

func main() {
	c = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	ctx := context.Background()
	go cacheValues(ctx)
	go snapshot(ctx)

	http.HandleFunc("/values", valuesHandler)
	http.ListenAndServe("127.0.0.1:8888", nil)
}
