package main

import (
	"fmt"
	"reflect"
)

type Sa struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type Sv struct {
	Key   string
	Value string
}

func main() {
	// v := "b"
	sa := Sa{
		A: 1,
		B: "b",
	}

	// use Key: B as sa's Field B, and change sa.B to sv.Value
	sv := Sv{
		Key:   "B",
		Value: "c",
	}

	t := reflect.TypeOf(sa)
	fmt.Println(t)

	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Name, t.Field(i).Tag)
		if t.Field(i).Name == sv.Key {
			// set sa.B to v.Value
			v := reflect.ValueOf(&sa)
			e := v.Elem()
			f := e.FieldByName(sv.Key)
			if f.IsValid() {
				f.SetString(sv.Value)
			}
		}
	}

	fmt.Println(sa)

}
