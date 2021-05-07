package main

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/Augustu/go-draft/redis/queue/common"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// var (
// 	addTask0, addTask1, addTask2                      tasks.Signature
// 	multiplyTask0, multiplyTask1                      tasks.Signature
// 	sumIntsTask, sumFloatsTask, concatTask, splitTask tasks.Signature
// 	panicTask                                         tasks.Signature
// 	longRunningTask                                   tasks.Signature
// )

var (
	server *machinery.Server
	err    error
)

func init() {

	server, err = common.StartServer()
	if err != nil {
		fmt.Println("start server failed: ", err)
	}

}

func send() error {
	var addTask0 tasks.Signature

	var initTasks = func() {
		addTask0 = tasks.Signature{
			Name: "add",
			Args: []tasks.Arg{
				{
					Type:  "int64",
					Value: 1,
				},
				{
					Type:  "int64",
					Value: 1,
				},
			},
		}
	}

	initTasks()

	// fmt.Println("Single task:")

	ctx := context.Background()

	addTask0.OnSuccess = []*tasks.Signature{
		{
			Name: "success",
			Args: []tasks.Arg{},
		},
	}

	addTask0.OnError = []*tasks.Signature{
		{
			Name: "fail",
			Args: []tasks.Arg{},
		},
	}

	_, err = server.SendTaskWithContext(ctx, &addTask0)
	if err != nil {
		return fmt.Errorf("could not send task: %s", err.Error())
	}

	// results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	// if err != nil {
	// 	return fmt.Errorf("getting task result failed with error: %s", err.Error())
	// }
	// fmt.Printf("1 + 1 = %v\n", tasks.HumanReadableResults(results))

	// chain := tasks.NewChain()

	return nil
}

func main() {
	// err := send()
	// fmt.Println(err)

	ticker := time.NewTicker(time.Second)

	var count int64

	sendFunc := func() {
		for {
			err := send()
			if err != nil {
				fmt.Println(err)
				continue
			}
			atomic.AddInt64(&count, 1)
		}
	}

	thread := runtime.GOMAXPROCS(0) / 2
	fmt.Println(thread)
	// pool := make(chan struct{}, thread)
	for i := 0; i < thread; i++ {
		// pool <- struct{}{}
		go sendFunc()
	}
	// fmt.Println("set pool")

	// for range pool {
	// 	go sendFunc()
	// }

	for {
		select {
		case <-ticker.C:
			fmt.Printf("send %d in one second\n", atomic.LoadInt64(&count))
			atomic.StoreInt64(&count, 0)
		}
	}
}
