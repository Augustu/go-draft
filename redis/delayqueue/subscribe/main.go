package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Augustu/go-draft/redis/delayqueue/common"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func worker() error {
	consumerTag := "machinery_worker"

	server, err := common.StartServer()
	if err != nil {
		return err
	}

	var handled int64

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			fmt.Printf("handled %d\n", atomic.LoadInt64(&handled))
		}
	}()

	worker := server.NewWorker(consumerTag, 0)

	errorhandler := func(err error) {
		// fmt.Println("handler error", err)
	}

	pretaskhandler := func(signature *tasks.Signature) {
		// fmt.Println("start task handler for", signature.Name)
	}

	posttaskhandler := func(signature *tasks.Signature) {
		// fmt.Println("end task handler for", signature.Name)
		atomic.AddInt64(&handled, 1)
	}

	worker.SetPostTaskHandler(posttaskhandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(pretaskhandler)

	return worker.Launch()
}

func main() {
	worker()
}
