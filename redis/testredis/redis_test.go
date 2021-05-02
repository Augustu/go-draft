package testredis

import (
	"fmt"
	"testing"
)

// func TestRedis(t *testing.T) {
// 	err := Redis()
// 	if err != nil {
// 		t.Fail()
// 	}
// }

func TestRedisAuth(t *testing.T) {
	err := RedisAuth()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
