package main

import (
	"fmt"

	"github.com/Augustu/go-draft/redis/queue/common"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func worker() error {
	consumerTag := "machinery_worker"

	server, err := common.StartServer()
	if err != nil {
		return err
	}

	worker := server.NewWorker(consumerTag, 0)

	errorhandler := func(err error) {
		fmt.Println("handler error", err)
	}

	pretaskhandler := func(signature *tasks.Signature) {
		fmt.Println("start task handler for", signature.Name)
	}

	posttaskhandler := func(signature *tasks.Signature) {
		fmt.Println("end task handler for", signature.Name)
	}

	worker.SetPostTaskHandler(posttaskhandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(pretaskhandler)

	return worker.Launch()
}

func main() {
	// rc := redis.NewClient(&redis.Options{
	// 	Addr:       "127.0.0.1:6379",
	// 	DB:         0,
	// 	MaxRetries: 3,
	// })

	// rc.XAdd(rc.Context(), &redis.XAddArgs{
	// 	Stream:       "",
	// 	MaxLen:       100,
	// 	MaxLenApprox: 0,
	// 	ID:           "",
	// 	Values:       nil,
	// })
	worker()
}
