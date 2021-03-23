package copymap

import "testing"

func TestCopymap(t *testing.T) {
	o := &Options{}

	fields := map[string]interface{}{
		"a": "aa",
		"b": "bb",
	}

	o.Copymap(fields)

	o.PrintField()
}
