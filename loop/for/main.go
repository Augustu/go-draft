package main

import "fmt"

func main() {
	total := -3
	for i := 0; i < total; i++ {
		fmt.Println("i = ", i)
	}

	list := []string{"a", "b", "c"}
	for idx, l := range list {
		fmt.Println(idx, l)
		if l == "b" {
			fmt.Println("continue", idx, l)
			continue
		}
		fmt.Println(idx, l)
	}

}
