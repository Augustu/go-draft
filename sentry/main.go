package main

import (
	"fmt"
	"time"

	log "github.com/Augustu/go-micro/v2/logger"
	"github.com/getsentry/raven-go"
	"github.com/getsentry/sentry-go"
)

func main1() {
	raven.SetDSN("http://27af8e6e1b22460793bc405dcdb6a88c@134.175.142.254:9000/29")

	raven.CapturePanic(test, map[string]string{"a": "b"})

	raven.CaptureError(fmt.Errorf("test"), nil)
}

func test() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		// sentry.CaptureException(err)
	// 		fmt.Println(err)
	// 	}
	// }()

	panic("test")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			// sentry.CaptureException(err)
			fmt.Println(err)
		}
	}()

	sc := sentry.ClientOptions{
		Dsn:              "test",
		AttachStacktrace: true,
	}

	err := sentry.Init(sc)
	if err != nil {
		log.Warnf("init sentry log failed: %s, start without sentry now", err)
		return
	}

	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("sentry log started")
	sentry.CaptureMessage("sentry log test started")
	sentry.CaptureException(fmt.Errorf("test"))

	log.Infof("sentry started")
	log.Warnf("sentry warn")
	// log.Fatal("sentry fatal")
	panic("test")
	log.Warnf("sentry warn")
}
