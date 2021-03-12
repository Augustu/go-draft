package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	addr := "http://127.0.0.1:8000/testcoll/v1"
	resp, err := http.Get(addr)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
