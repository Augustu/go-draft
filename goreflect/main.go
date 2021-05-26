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

	zi := 0
	zs := ""

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		var a interface{}
		var z interface{}

		f := t.Field(i)

		k := v.Field(i).Kind()
		ty := v.Field(i).Type()

		fmt.Println(t.Field(i).Name, k, ty)

		if strings.Contains(ty.String(), "int") {
			n, ok := convertInt(u.Value)
			if ok {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &n
					z = &zi
				} else {
					a = n
					z = zi
				}
			}
		} else if strings.Contains(ty.String(), "string") {
			n, ok := convertString(u.Value)
			if ok {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &n
					z = &zs
				} else {
					a = n
					z = zs
				}
			}
		}

		fmt.Println("za", z, a, k.String(), ty.String())

		if f.Tag.Get("json") == u.Key {
			v.Field(i).Set(reflect.ValueOf(a))
		} else {
			v.Field(i).Set(reflect.ValueOf(z))
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

func main2() {
	var a int = 1

	fmt.Println(reflect.TypeOf(a))
	fmt.Println(reflect.TypeOf(&a))

	t := reflect.TypeOf(&a)
	fmt.Println(t.Kind(), t.Elem(), t.Elem().Kind())
}
