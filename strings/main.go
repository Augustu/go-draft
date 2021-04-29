package main

import (
	"fmt"
	"strings"
)

func mainCase() {
	a := "test"

	at := strings.ToLower(a)
	fmt.Println(at)

	// at[0] = strings.ToUpper(string(at[0]))[0:1]

	h := at[0]
	hs := strings.ToUpper(string(h))

	fmt.Println(at)

	fmt.Printf("%s%s", hs, at[1:])
}

func main() {
	test := "a-b-c--d"
	// test := "-a-b"
	l := strings.Split(test, "-")
	fmt.Println(len(l), l)
}
