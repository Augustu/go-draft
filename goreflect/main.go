package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type A struct {
	A int     `json:"a"`
	B string  `json:"b"`
	C *string `json:"c"`
}

type B struct {
	A string  `json:"a"`
	B int     `json:"b"`
	C *int    `json:"c"`
	D *string `json:"d"`
}

type U struct {
	Key   string
	Value string
}

type V struct {
	Key   string
	Value interface{}
}

func updateStruct(s interface{}, u U) {

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f.Type, f.Name, f.Tag, v.Field(i))
		if f.Tag.Get("json") == u.Key {
			var a interface{}
			if v.Field(i).Kind() == reflect.Ptr {
				a = &u.Value
			} else {
				a = u.Value
			}
			v.Field(i).Set(reflect.ValueOf(a))
		}
	}

}

// notice here u is V struct
func updateStructInterface(s interface{}, u V) {

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)

		if f.Tag.Get("json") == u.Key {

			s := v.Field(i).Type().String()
			if strings.Contains(s, "int") {
				n, ok := convertInt(u.Value)
				if ok {
					var a interface{}
					if v.Field(i).Kind() == reflect.Ptr {
						a = &n
					} else {
						a = n
					}
					v.Field(i).Set(reflect.ValueOf(a))
				}

			} else if strings.Contains(s, "string") {
				n, ok := convertString(u.Value)
				if ok {
					var a interface{}
					if v.Field(i).Kind() == reflect.Ptr {
						a = &n
					} else {
						a = n
					}
					v.Field(i).Set(reflect.ValueOf(a))
				}
			}

		}

	}

}

func convertString(v interface{}) (string, bool) {
	s, ok := v.(string)
	if ok {
		return s, true
	}

	n, ok := v.(int)
	if ok {
		return fmt.Sprint(n), true
	}

	return "", false
}

func convertInt(v interface{}) (int, bool) {
	i, ok := v.(int)
	if ok {
		return i, true
	}

	s, ok := v.(string)
	if ok {
		i, e := strconv.Atoi(s)
		if e == nil {
			return i, true
		}
	}

	return 0, false
}

func main() {
	/*
		a := A{
			A: 1,
			B: "2",
		}

			u := U{
				Key:   "b",
				Value: "3",
			}

			uu := U{
				Key:   "c",
				Value: "3",
			}

			updateStruct(&a, u)
			updateStruct(&a, uu)

			fmt.Println(a, *a.C)

			// updateStructInterface(&a, v)
			// fmt.Println(a, *a.C)

	*/

	v := V{
		Key:   "b",
		Value: "5",
	}

	b := B{
		A: "1",
		B: 2,
	}

	updateStructInterface(&b, v)
	fmt.Println(b, b.C)

	v.Key = "c"
	updateStructInterface(&b, v)
	fmt.Println(b, *b.C)

	v.Key = "a"
	updateStructInterface(&b, v)
	fmt.Println(b, *b.C)

	v.Key = "d"
	updateStructInterface(&b, v)
	fmt.Println(b, *b.D)
}
