package singleinstance

import (
	"testing"
)

func TestLock(t *testing.T) {
	// key := "testkey"

	// res := "resource"

	// const total int = 500

	// ch := make(chan bool)

	// for i := 0; i < total; i++ {
	// 	go func() {
	// 		for i := 0; i < total; i++ {
	// 			Lock(key)
	// 			_, err := lock.client.Incr(lock.context, res).Result()
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
