package main

import "fmt"

type Options struct {
	Fields map[string]interface{}
}

func (o *Options) Copymap(fields map[string]interface{}) {
	o.Fields = fields
	// o.Fields = copyFields(fields)
	fmt.Printf("Addr o: %p, field: %p\n", &(o.Fields), &fields)
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (o *Options) PrintField() {
	for k, v := range o.Fields {
		fmt.Println(k, v)
	}
}

func main() {
	o := &Options{}

	fields := map[string]interface{}{
		"a": "aa",
		"b": "bb",
	}

	o.Copymap(fields)

	o.PrintField()

}
