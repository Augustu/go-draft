package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Augustu/go-micro/v2/logger"
	"github.com/getsentry/sentry-go"
)

const (
	serverAddr string = "127.0.0.1:8888"
)

func initSentry() {
	sc := sentry.ClientOptions{
		Dsn: "test",
	}

	err := sentry.Init(sc)
	if err != nil {
		log.Warnf("init sentry log failed: %s, start without sentry now", err)
		return
	}

	defer sentry.Flush(2 * time.Second)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("recovering with error:  ", err)
			sentry.CaptureException(fmt.Errorf("index handler failed: %s", err))
			return
		}
	}()

	fmt.Fprint(w, "index")
	panic("index panic")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover http server")
			panic(http.ListenAndServe(serverAddr, nil))
		}
	}()

	initSentry()

	http.HandleFunc("/index", indexHandler)

	panic(http.ListenAndServe(serverAddr, nil))
}
