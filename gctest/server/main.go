package main

import (
	"net/http"
)

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("a"))

	// if r != nil {
	// 	r.Body.Close()
	// }

	// r.Body.Close()
}

func main() {
	// go func() {
	// 	for {
	// 		runtime.GC()
	// 		time.Sleep(time.Millisecond)
	// 	}
	// }()

	http.HandleFunc("/test", test)
	http.ListenAndServe("127.0.0.1:8000", nil)
}
