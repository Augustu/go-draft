package main

import (
	"fmt"
	"time"

	log "github.com/Augustu/go-micro/v2/logger"
	"github.com/getsentry/sentry-go"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			// sentry.CaptureException(err)
			fmt.Println(err)
		}
	}()

	sc := sentry.ClientOptions{
		Dsn: "test",
	}

	err := sentry.Init(sc)
	if err != nil {
		log.Warnf("init sentry log failed: %s, start without sentry now", err)
		return
	}

	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("sentry log started")
	sentry.CaptureMessage("sentry log test started")

	log.Infof("sentry started")
	log.Warnf("sentry warn")
	// log.Fatal("sentry fatal")
	panic("test")
	log.Warnf("sentry warn")
}
