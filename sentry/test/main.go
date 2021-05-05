package main

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	dsn string = "http://48078ee986324a05b2c30741203e0b59@localhost:8080/2"
)

func testPanic() {
	defer func() {
		r := recover()
		sentry.Logger.Println("logger")
		sentry.CaptureException(fmt.Errorf("%s", r))
		fmt.Println("recovered")
	}()
	// panic("test")
	var t *[]string
	for range *t {
	}
}

func main() {
	sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
	})

	// sentry.CaptureException(errors.New("test 1 error"))

	testPanic()

	sentry.Flush(time.Second * 1)

	// panic("test")
}
