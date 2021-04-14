package testredis

import (
	"fmt"
	"log"
	"time"

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

func RedisAuth() error {
	rc := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 3,
		Password:   "test",
	})

	// res, err := rc.Pipeline().Auth("testpass").Result()
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(res)

	// res, err := rc.Do("AUTH", "test").Result()
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(res)

	s, err := rc.Ping().Result()
	if err != nil {
		return err
	}

	_, err = rc.Set("test", "test", time.Second).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
	return nil
}
