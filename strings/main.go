package main

import (
	"fmt"
	"sort"
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

func mainSplit() {
	test := "a-b-c--d"
	// test := "-a-b"
	l := strings.Split(test, "-")
	fmt.Println(len(l), l)
}

func main() {
	cases := []string{
		"a-20210301",
		"a-20210302",
		"a-20210303",
	}

	sort.Strings(cases)

	fmt.Println(cases)

	in := []int{
		20210203,
		20210204,
		20210130,
	}

	sort.Ints(in)

	fmt.Println(in)
}
