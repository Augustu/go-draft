package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type client struct {
	*redis.Client
	context context.Context
}

type sa struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func (c *client) put(key string, total int) {
	for i := 0; i < total; i++ {

		b, err := json.Marshal(sa{
			A: i,
			B: strconv.Itoa(i),
		})
		if err != nil {
			log.Println(err)
		}

		_, err = c.RPush(c.context, key, string(b)).Result()
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *client) get(key string) {
	for {
		var ss sa
		res, err := c.BLPop(c.context, time.Second, key).Result()
		if err != nil {
			log.Println(err)
			break
		}

		s := res[1]
		err = json.Unmarshal([]byte(s), &ss)
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println(ss)
	}
}

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         1,
		MaxRetries: 3,
	})

	c := &client{
		context: context.Background(),
		Client:  rc,
	}

	// go c.put("a", 10000)

	c.get("a")
	// select {}
}
