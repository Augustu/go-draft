package testredis

import (
	"fmt"

	"github.com/go-redis/redis"
)

func Redis() error {
	rc := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 3,
	})

	s, err := rc.Ping().Result()
	if err != nil {
		return err
	}

	fmt.Println(s)
	return nil
}
