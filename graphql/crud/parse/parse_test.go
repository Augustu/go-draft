package parse

import (
	"fmt"
	"testing"
)

var (
	// {list(info: "123", where: [{price: {lt: 1}}, {id: {gt: 1, lt: 5}, price: {gt: 1.33}}]){id,name,info,price}}
	t1 = []interface{}{
		map[string]interface{}{
			"price": map[string]interface{}{
				"lt": 1,
			},
		},
		map[string]interface{}{
			"id": map[string]interface{}{
				"gt": 1,
				"lt": 5,
			},
			"price": map[string]interface{}{
				"gt": 1.33,
			},
		},
	}
)

func TestQuery(t *testing.T) {
	r := Query(t1)
	fmt.Printf("query: %#v\n", r)
}
