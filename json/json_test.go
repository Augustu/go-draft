package json

import (
	"bytes"
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
		B: "b<b",
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

	b, err := json.MarshalIndent(tt, "", "    ")
	if err != nil {
		t.Fail()
	}

	fmt.Println("json indent", string(b))
	fmt.Printf("format %s\n", fmt.Sprintf("%s", b))

	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(tt)
	if err != nil {
		t.Fail()
	}
	fmt.Println("Escape html", buffer.String())

	var ttt test
	if err = json.Unmarshal(body, &ttt); err != nil {
		t.Fail()
	}
	fmt.Println(ttt)
}
