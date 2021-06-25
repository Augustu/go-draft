package main

import (
	"fmt"
	"sync"
)

func main() {
	sm := sync.Map{}

	sm.Store("a", "b")

	v, ok := sm.Load("a")
	fmt.Println(v, ok)

	sm.Range(func(key, value interface{}) bool {
		sm.Delete(key)
		return true
	})

	sm.Range(func(key, value interface{}) bool {
		fmt.Println(key)
		return true
	})
}
