package testredis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// func Redis() error {
// 	rc := redis.NewClient(&redis.Options{
// 		Addr:       "127.0.0.1:6379",
// 		DB:         0,
// 		MaxRetries: 3,
// 	})

// 	s, err := rc.Ping().Result()
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(s)
// 	return nil
// }

func RedisPipeline() error {
	rc := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 3,
		// Password:   "test",
	})

	p := rc.Pipeline()
	ic := p.HSet(rc.Context(), "a", "b", "c")
	fmt.Println(ic)

	c, err := p.Exec(rc.Context())
	if err != nil {
		fmt.Println(err)
		return err
	}

	for idx, cmd := range c {
		fmt.Println(idx, cmd.Args(), cmd.String(), cmd.Err())
	}
	fmt.Println(c)

	return nil
}

func RedisAuth() error {
	rc := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 3,
		// Password:   "test",
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
	ctx := context.Background()

	s, err := rc.Ping(ctx).Result()
	if err != nil {
		return err
	}

	_, err = rc.Set(ctx, "test", "test", time.Second).Result()
	if err != nil {
		log.Fatal(err)
	}

	_, err = rc.Set(ctx, "test", "test", 0).Result()
	if err != nil {
		log.Fatal(err)
	}

	_, err = rc.Set(ctx, "test", "ttt", 0).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
	return nil
}
