package parseurl_test

import (
	"testing"

	"github.com/Augustu/go-draft/parseurl"
)

func TestParseURL(t *testing.T) {
	tests := []string{
		"http://www.google.com/?q=go+language#foo%26bar",
		"http://www.google.com/file%20one%26two",
	}

	for _, t := range tests {
		parseurl.ParseURL(t)
	}
}
