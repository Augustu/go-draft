package main

import (
	"fmt"
	"net/http"
)

func handler(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	fmt.Println(req.Header)

	values := req.URL.Query()
	fmt.Println("Values: ", values.Encode())
	fmt.Println("Query a = ", values.Get("a"))

	fmt.Println("a = ", req.Header.Get("a"))
	fmt.Println("handle done")

	resp.Write([]byte("done"))
}

func main() {
	http.HandleFunc("/api/v*", handler)

	http.ListenAndServe("127.0.0.1:8080", nil)
}
