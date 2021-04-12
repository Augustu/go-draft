package singleinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ttl = 3 * time.Second
)

type locker struct {
	client  *redis.Client
	timer   *time.Ticker
	context context.Context
	cancel  context.CancelFunc
}

func (l *locker) Lock(key string) error {
	res, err := l.client.SetNX(l.context, key, true, ttl).Result()
	if err != nil {
		return fmt.Errorf("lock key %s failed: %s", key, err.Error())
	}

	if !res {
		return fmt.Errorf("key %s already locked", key)
	}

	// feed lock with ttl
	go func() {
		l.timer = time.NewTicker(ttl / 2)

		for {
			select {
			case <-l.context.Done():
				if err = l.Unlock(key); err != nil {
					fmt.Printf("main thread exit, but unlock key %s failed: %s", key, err.Error())
				}

			case <-l.timer.C:
				if err = l.feed(key); err != nil {
					fmt.Printf("feed key %s failed: %s", key, err.Error())
				}
			}
		}
	}()

	return nil
}

func (l *locker) feed(key string) error {

	retry := 3
	i := 0

	// var res string
	var err error

	for i = 0; i < retry; i++ {
		_, err = l.client.SetEX(l.context, key, true, ttl).Result()
		if err != nil {
			continue
		}

		// fmt.Printf("debug res: %s\n", res)
		return nil
	}

	if i == retry {
		return fmt.Errorf("feed key %s failed: %s", key, err.Error())
	}
	return nil
}

func (l *locker) Unlock(key string) error {

	_, err := l.client.Del(l.context, key).Result()
	if err != nil {
		return fmt.Errorf("unlock key %s failed: %s", key, err.Error())
	}
	l.timer.Stop()

	// fmt.Printf("debug res: %d\n", res)
	return nil
}

func (l *locker) Close() error {
	l.cancel()
	return l.client.Close()
}

var lock locker

func init() {
	c := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         1,
		MaxRetries: 3,
	})

	ctx, cancel := context.WithCancel(context.Background())

	lock = locker{
		client:  c,
		context: ctx,
		cancel:  cancel,
	}
}

func Lock(key string) {
	for {
		err := lock.Lock(key)
		if err != nil {
			return
		}
		time.Sleep(time.Microsecond)
	}
}

func Unlock(key string) error {
	return lock.Unlock(key)
}

func Close() error {
	return lock.Close()
}

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         1,
		MaxRetries: 3,
	})
}
