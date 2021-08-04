package main

import "fmt"

func main() {
	a := make([]int, 3)
	fmt.Printf("%p\n", &a)

	a = a[0:0]
	fmt.Printf("%p\n", &a)

	b := make([]int, 5)
	fmt.Printf("%p\n", &b)

	cleanSlice(&b)
	fmt.Printf("%p\n", &b)

	cleanSlice2(&b)
	fmt.Printf("%p\n", &b)
}

func cleanSlice(s *[]int) {
	// *s = append([]int{})
	*s = []int{}
}

func cleanSlice2(s *[]int) {
	*s = (*s)[0:0]
}
