package singleinstance

import (
	"testing"
	"time"
)

func TestLocker(t *testing.T) {
	var err error

	key := "testkey"

	err = Lock(key)
	if err != nil {
		t.Fail()
	}

	// try lock twice
	err = Lock(key)
	if err == nil {
		t.Fail()
	}

	time.Sleep(5 * time.Second)

	// try another lock
	err = Lock(key)
	if err == nil {
		t.Fail()
	}

	err = Unlock(key)
	if err != nil {
		t.Fail()
	}
}
