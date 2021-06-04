package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Augustu/go-draft/redis/delayqueue/common"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
)

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

		eta := time.Now().Add(5 * time.Second)
		fmt.Println("eta: ", eta.String())
		addTask0.ETA = &eta
	}

	initTasks()

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

	return nil
}

type A struct {
	A string `json:"a"`
	B string `json:"b"`
}

func sendNotification() error {

	eta := time.Now().Add(3 * time.Second)
	fmt.Println(eta)

	notiTask := tasks.Signature{
		Name: "sendNotification",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: "notify body string " + eta.String(),
				// struct can not work
				// Value: A{
				// 	A: "a",
				// 	B: "b",
				// },
			},
		},
	}

	notiTask.ETA = &eta

	ctx := context.Background()

	_, err = server.SendTaskWithContext(ctx, &notiTask)
	if err != nil {
		return fmt.Errorf("could not send task: %s", err.Error())
	}

	return nil
}

func main() {
	ticker := time.NewTicker(time.Second)

	var need int64 = 10

	var total int64
	var count int64

	mux := sync.Mutex{}

	sendFunc := func() {
		for {

			// lock
			mux.Lock()

			// check before send
			if atomic.CompareAndSwapInt64(&total, need, need) {
				return
			}

			err := sendNotification()
			if err != nil {
				fmt.Println(err)
				continue
			}

			atomic.AddInt64(&count, 1)
			atomic.AddInt64(&total, 1)

			// unlock
			mux.Unlock()
		}
	}

	thread := runtime.GOMAXPROCS(0) / 2
	fmt.Println("total thread", thread)

	for i := 0; i < thread; i++ {
		go sendFunc()
	}

	for range ticker.C {
		fmt.Printf("send %d in one second, total %d\n", atomic.LoadInt64(&count), atomic.LoadInt64(&total))
		atomic.StoreInt64(&count, 0)
	}
}
