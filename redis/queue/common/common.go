package common

import (
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

func StartServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Broker:          "redis://localhost:6379/0",
		ResultBackend:   "redis://localhost:6379/1",
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	// Register tasks
	tasks := map[string]interface{}{
		"add":     Add,
		"success": Success,
		"fail":    Fail,
	}

	return server, server.RegisterTasks(tasks)
}

func Add(args ...int64) (int64, error) {
	fmt.Println("add", time.Now().String())
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}

	// time.Sleep(30 * time.Second)

	return sum, nil
}

func Success(res int64) (string, error) {
	fmt.Println("success", time.Now().String())
	return fmt.Sprintf("done now: %d", res), nil
}

func Fail(msg string) (string, error) {
	fmt.Println("failed", time.Now().String())
	return msg, nil
}
