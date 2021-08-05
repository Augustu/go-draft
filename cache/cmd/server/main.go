package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Augustu/go-draft/cache/cmd/server/model"
	"github.com/Augustu/go-draft/cache/pkg/config"
	"github.com/Augustu/go-draft/cache/pkg/store"
)

const (
	dsn string = "root:FuCkU@!@#$%^@tcp(10.10.15.11:32306)/cache?charset=utf8mb4&parseTime=True&loc=Local"

	host   string = "10.10.15.11:32379"
	passwd string = ""

	addr string = "127.0.0.1:8000"
)

func newConfig() config.Config {
	c := config.Config{}
	c.URI = dsn
	c.Host = host
	c.Passwd = passwd
	return c
}

var (
	statsChan = make(chan model.Stats, 10240)
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	s := model.Stats{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return
	}

	statsChan <- s
}

func storeHandler(c *config.Store) {
	st := store.New(c)

	st.DB.AutoMigrate(&model.Stats{})

	batchsize := 1000
	ticker := time.NewTicker(time.Second)

	ss := []model.Stats{}

	for {
		select {
		case s := <-statsChan:
			fmt.Println("got stats", s)
			ss = append(ss, s)
			if len(ss) >= batchsize {
				writeStore(st, ss, statsChan)
				ss = []model.Stats{}
			}

		case <-ticker.C:
			if len(ss) > 0 {
				writeStore(st, ss, statsChan)
				ss = []model.Stats{}
			}
		}
	}

}

func writeStore(st *store.Store, ss []model.Stats, fallback chan<- model.Stats) {
	err := st.DB.CreateInBatches(&ss, 1000).Error
	if err != nil {
		for _, s := range ss {
			fallback <- s
		}
	}
}

func main() {

	stopChan := make(chan struct{}, 1)

	go func() {
		select {
		case <-stopChan:
			return

		default:
			http.HandleFunc("/create", statsHandler)
			http.ListenAndServe(addr, nil)
		}
	}()

	cf := newConfig()

	for i := 0; i < 8; i++ {
		go func(s *config.Store) {
			storeHandler(s)
		}(&cf.Store)
	}

	shutdownHandler := make(chan os.Signal, 2)
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

	signal.Notify(shutdownHandler, shutdownSignals...)

	fmt.Println("listen and serve, waiting")
	<-shutdownHandler
	close(stopChan)
	close(statsChan)
}
