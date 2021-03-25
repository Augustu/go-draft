package maps

import "fmt"

type root struct {
	m map[string]string
}

type childa struct {
	m *map[string]string
}

type childb struct {
	m *map[string]string
}

type childc struct {
	m *map[string]string
}

func ChildMap() {
	r := root{
		m: make(map[string]string),
	}

	a := childa{
		m: &r.m,
	}
	b := childb{
		m: &r.m,
	}
	c := childc{
		m: &r.m,
	}

	r.m["a"] = "aa"
	fmt.Println(r.m["a"])
	fmt.Println((*a.m)["a"])
	fmt.Println((*b.m)["a"])
	fmt.Println((*c.m)["a"])
}
