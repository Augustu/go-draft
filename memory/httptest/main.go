package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/Augustu/go-draft/memory/httptest/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsn string = "root:123456@tcp(192.168.0.10:33066)/test?charset=utf8mb4&parseTime=True&loc=Local"
)

type client struct {
	db *gorm.DB

	req chan model.User
}

func new() *client {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Printf("create mysql db with dsn: %s failed: %s", dsn, err)
		return nil
	}

	db.AutoMigrate(&model.User{})

	return &client{
		db:  db,
		req: make(chan model.User, 1024000),
	}
}

func (c *client) handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var u model.User
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.req <- u
}

func (c *client) regist() {
	for u := range c.req {
		go func(u model.User) {
			err := c.db.Create(&u).Error
			if err != nil {
				log.Printf("create user: %#v failed: %s", u, err)
			}
		}(u)
	}
}

func (c *client) noop() {
	for range c.req {
		go func() {
			time.Sleep(time.Second)
		}()
	}
}

func (c *client) monitor() {
	ticker := time.NewTicker(time.Millisecond)

	for range ticker.C {
		// runtime.GC()
		log.Printf("num of go runtines: %d", runtime.NumGoroutine())
	}
}

func main() {
	c := new()
	// go c.regist()
	go c.monitor()
	go c.noop()

	http.HandleFunc("/test", c.handler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}
