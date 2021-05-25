package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type client struct {
	*redis.Client
	context context.Context
}

func (c *client) put(key string, total int) {
	for i := 0; i < total; i++ {
		c.RPush(c.context, key, i)
	}
}

func (c *client) get(key string) {
	for {
		res, err := c.BLPop(c.context, time.Second, key).Result()
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(res)
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

	go c.put("a", 10)

	c.get("a")
}
