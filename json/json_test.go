package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

type test struct {
	a string `json:"a"`
	B string `json:"b"`
}

func TestMarshal(t *testing.T) {
	tt := test{
		a: "aa",
		B: "bb",
	}

	body, err := json.Marshal(tt)
	if err != nil {
		t.Fail()
	}

	fmt.Println(string(body))
}
