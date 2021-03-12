package parseurl

import (
	"fmt"
	"net/url"
)

func ParseURL(query string) {
	// url.ParseQuery(url)
	values, err := url.ParseQuery(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Values: %#v\n", values)

	u, err := url.ParseRequestURI(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Url: %v", u)
}
