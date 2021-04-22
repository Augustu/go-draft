package json

import (
	"fmt"
	"testing"

	json "github.com/json-iterator/go"
)

type test struct {
	a   string            `json:"a"`
	B   string            `json:"b"`
	Map map[string]string `json:"map"`
}

func TestMarshal(t *testing.T) {
	tt := test{
		a: "aa",
		B: "bb",
		Map: map[string]string{
			"a": "aaa",
		},
	}

	body, err := json.Marshal(tt)
	if err != nil {
		t.Fail()
	}

	res := json.Get(body, "map", "a")
	fmt.Println(res.ToString())

	fmt.Println(string(body))

	var ttt test
	if err = json.Unmarshal(body, &ttt); err != nil {
		t.Fail()
	}
	fmt.Println(ttt)
}
