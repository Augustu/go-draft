package singleinstance

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestLock(t *testing.T) {
	// key := "testkey"

	// res := "resource"

	// const total int = 500

	// ch := make(chan bool)

	// for i := 0; i < total; i++ {
	// 	go func() {
	// 		c := NewClient()
	// 		for i := 0; i < total; i++ {
	// 			Lock(key)
	// 			_, err := c.Incr(lock.context, res).Result()
	// 			if err != nil {
	// 				t.Fail()
	// 			}
	// 			Unlock(key)
	// 		}
	// 		ch <- true
	// 	}()
	// }

	// for i := 0; i < total; i++ {
	// 	<-ch
	// }

}

func TestDup(t *testing.T) {
	lockKey := "LOCKER"
	key := "testkey"

	const total int = 500

	ch := make(chan bool)

	for i := 0; i < total; i++ {
		go func() {
			c := NewClient()
			for i := 0; i < total; i++ {
				_, err := c.ZAdd(lock.context, key, &redis.Z{
					Score:  rand.Float64(),
					Member: rand.Int31(),
				}).Result()
				if err != nil {
					log.Print(err)
					t.Fail()
				}
			}
			ch <- true
		}()
	}

	c := NewClient()
	for i := 0; i < 5; i++ {
		<-ch
		// _, err := c.ZUnionStore(lock.context, fmt.Sprintf("%s-%d", key, i), &redis.ZStore{Keys: []string{key}}).Result()
		// if err != nil {
		// 	t.Fail()
		// }
		// _, err = c.Del(lock.context, key).Result()
		// if err != nil {
		// 	t.Fail()
		// }

		Lock(lockKey)
		_, err := c.RenameNX(lock.context, key, fmt.Sprintf("%s-%d", key, i)).Result()
		if err != nil {
			log.Print(err)
			t.Fail()
		}
		Unlock(lockKey)
	}

}

func TestAdd(t *testing.T) {
	lockKey := "LOCKER"
	key := "testkey"

	const total int = 500

	ch := make(chan bool)

	for i := 0; i < total; i++ {
		go func() {
			c := NewClient()
			for i := 0; i < total; i++ {
				// Lock(lockKey)
				_, err := c.ZIncrBy(lock.context, key, 1, "m").Result()
				if err != nil {
					t.Fail()
				}
				// Unlock(lockKey)
			}
			ch <- true
		}()
	}

	c := NewClient()
	var h []string

	for i := 0; i < total; i++ {
		<-ch
		if i%11 == 0 {
			Lock(lockKey)

			nk := fmt.Sprintf("%s-%d", key, i)

			c.Rename(lock.context, key, nk).Result()

			// _, err := c.ZUnionStore(lock.context, nk, &redis.ZStore{Keys: []string{key}}).Result()
			// if err != nil {
			// 	t.Fail()
			// }
			// _, err = c.Del(lock.context, key).Result()
			// if err != nil {
			// 	t.Fail()
			// }
			Unlock(lockKey)

			h = append(h, nk)
		}
	}

	var count float64 = 0

	for _, s := range h {
		score, _ := c.ZScore(lock.context, s, "m").Result()

		count += score
	}

	fmt.Printf("count: %f\n", count)

	if count != 250000 {
		t.Fail()
	}

	// c := NewClient()
	// for i := 0; i < 5; i++ {
	// 	<-ch
	// 	// _, err := c.ZUnionStore(lock.context, fmt.Sprintf("%s-%d", key, i), &redis.ZStore{Keys: []string{key}}).Result()
	// 	// if err != nil {
	// 	// 	t.Fail()
	// 	// }
	// 	// _, err = c.Del(lock.context, key).Result()
	// 	// if err != nil {
	// 	// 	t.Fail()
	// 	// }

	// 	Lock(lockKey)
	// 	_, err := c.RenameNX(lock.context, key, fmt.Sprintf("%s-%d", key, i)).Result()
	// 	if err != nil {
	// 		log.Print(err)
	// 		t.Fail()
	// 	}
	// 	Unlock(lockKey)
	// }

}

func TestLocker(t *testing.T) {
	// TODO fix ci redis, or find another way to test

	// var err error

	// key := "testkey"

	// err = Lock(key)
	// if err != nil {
	// 	t.Fail()
	// }

	// // try lock twice
	// err = Lock(key)
	// if err == nil {
	// 	t.Fail()
	// }

	// time.Sleep(5 * time.Second)

	// // try another lock
	// err = Lock(key)
	// if err == nil {
	// 	t.Fail()
	// }

	// err = Unlock(key)
	// if err != nil {
	// 	t.Fail()
	// }
}
