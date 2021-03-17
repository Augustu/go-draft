package main

import (
	"fmt"
	"log"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
)

func main() {
	// etcSource := etcd.NewSource(
	// 	etcd.WithAddress("127.0.0.1:2379"),
	// 	etcd.WithPrefix("/mirco/config"),
	// 	// etcd.StripPrefix(true),
	// )

	// conf, _ := config.NewConfig()
	// conf.Load(etcSource)

	fileSource := file.NewSource()

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("NewConfig: %s", err.Error())
	}
	err = conf.Load(fileSource)
	if err != nil {
		log.Fatalf("Load file: %s", err.Error())
	}

	// watcher, _ := conf.Watch("mirco", "config", "database")
	watcher, err := conf.Watch("database")
	if err != nil {
		log.Fatalf("Watch database: %s", err.Error())
	}

	// envSource := env.NewSource()

	// conf, err := config.NewConfig()
	// if err != nil {
	// 	log.Fatalf("NewConfig: %s", err.Error())
	// }
	// err = conf.Load(envSource)
	// if err != nil {
	// 	log.Fatalf("Load file: %s", err.Error())
	// }

	// // watcher, _ := conf.Watch("mirco", "config", "database")
	// watcher, err := conf.Watch("database")
	// if err != nil {
	// 	log.Fatalf("Watch database: %s", err.Error())
	// }

	for {
		value, err := watcher.Next()
		if err != nil {
			log.Fatalf("Watcher next: %s", err.Error())
		}

		fmt.Println(string(value.Bytes()))
		fmt.Println(value.String("database"))
	}
}
