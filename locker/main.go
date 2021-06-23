package main

import "sync"

func main() {
	rwmux := sync.RWMutex{}

	rwmux.RLock()

	rwmux.RUnlock()
	rwmux.RUnlock()
	// rwmux.Unlock()
	// rwmux.Unlock()

	// mux := sync.Mutex{}

	// mux.Lock()

	// mux.Unlock()
	// mux.Unlock()
}
