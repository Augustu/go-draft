package testredis

import "testing"

func TestRedis(t *testing.T) {
	err := Redis()
	if err != nil {
		t.Fail()
	}
}
