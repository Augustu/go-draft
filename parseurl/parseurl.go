package parseurl

import (
	"net/url"
)

func ParseURL(query string) {
	// url.ParseQuery(url)
	_, err := url.ParseQuery(query)
	if err != nil {
		panic(err)
	}

	_, err = url.ParseRequestURI(query)
	if err != nil {
		panic(err)
	}
}
