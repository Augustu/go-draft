package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Augustu/go-draft/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataReport struct {
	gorm.Model
	App        string      `json:"app"`
	UID        string      `json:"uid"`
	EventID    string      `json:"event_id"`
	HappenedAt string      `json:"happened_at"`
	Attr       string      `json:"attr"`
	Param      string      `json:"param"`
	Value      interface{} `json:"value"`
}

type Recharge struct {
	CreatedAt      *time.Time `gorm:"index" json:"created_at,omitempty"`
	App            *string    `gorm:"type:varchar(16)" json:"app"`
	UID            *string    `gorm:"type:varchar(32)" json:"uid"`
	EventID        *string    `gorm:"type:varchar(64)" json:"event_id"`
	RechargeAmount *int64     `gorm:"type:bigint(11)" json:"recharge_amount"`
	HappenedAt     *time.Time `json:"happened_at"`
}

const (
	dsn string = "root:123456@tcp(127.0.0.1:3306)/mission?charset=utf8mb4\u0026parseTime=True\u0026loc=Local"
)

type client struct {
	db *gorm.DB
}

func NewClient() *client {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open mysql failed: %s", err)
	}

	return &client{
		db: db,
	}
}

func (c *client) DB() *gorm.DB {
	return c.db
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
}

func (c *client) rechargeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("read request failed: %s", err)
		return
	}

	dr := DataReport{}
	err = json.Unmarshal(body, &dr)
	if err != nil {
		log.Fatalf("unmarshal recharge failed: %s, body: %s", err, string(body))
		return
	}

	re, err := drToRecharge(&dr)
	if err != nil {
		log.Printf("convert data report to recharge failed: %s", err)
		return
	}

	// db := c.db.WithContext(context.Background())
	err = c.db.Create(&re).Error
	if err != nil {
		log.Fatalf("create recharge failed: %s", err)
		return
	}

}

func drToRecharge(dr *DataReport) (*Recharge, error) {
	happenedAt, err := utils.ParseTime(dr.HappenedAt)
	if err != nil {
		err = fmt.Errorf("parse happened at: %s failed: %s", dr.HappenedAt, err.Error())
		return nil, err
	}

	m := Recharge{
		App:     &dr.App,
		UID:     &dr.UID,
		EventID: &dr.EventID,
	}

	err = updateFields(&m, dr)
	if err != nil {
		err = fmt.Errorf("nothing to update: %#v", dr)
		return nil, err
	}

	m.App = &dr.App
	m.UID = &dr.UID
	m.EventID = &dr.EventID
	m.HappenedAt = &happenedAt

	return &m, nil
}

var memprofile = flag.String("memprofile", "", "write mem profile to file")

func main() {
	flag.Parse()

	var err error
	var memfile *os.File

	if *memprofile != "" {
		memfile, err = os.Create(*memprofile)
		if err != nil {
			log.Printf("create mem profile failed: %s", err.Error())
		}
	}

	c := NewClient()

	ctx, cancel := context.WithCancel(context.Background())

	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				// runtime.GC()
				ms := &runtime.MemStats{}
				runtime.ReadMemStats(ms)
				log.Printf("gc num: %d, HeapReleased: %d HeapInuse: %d, HeapIdle: %d, HeapAlloc: %d HeapSys: %d", ms.NumGC, ms.HeapReleased, ms.HeapInuse, ms.HeapIdle, ms.HeapAlloc, ms.HeapSys)
				log.Printf("gc num routines: %d", runtime.NumGoroutine())

			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		http.HandleFunc("/api/v1/data-report", c.rechargeHandler)
		// http.HandleFunc("/api/v1/data-report", handler)
		http.ListenAndServe("127.0.0.1:8000", nil)
		log.Print("start server at: 127.0.0.1:8000")

		<-ctx.Done()
		log.Print("exit server listener")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-sc
	cancel()

	if *memprofile != "" {
		pprof.WriteHeapProfile(memfile)
		memfile.Close()
	}
}

// updateFields init m 's time value and others before put in this
func updateFields(m interface{}, d *DataReport) error {
	zi := int(0)
	zs := ""
	zi64 := int64(0)

	v := reflect.ValueOf(m).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		var a interface{}
		var z interface{}

		f := t.Field(i)

		ty := v.Field(i).Type()

		if strings.Contains(ty.String(), "int64") {
			n, ok := convertInt64(d.Value)
			if ok {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &n
					z = &zi64
				} else {
					a = n
					z = zi64
				}
			} else {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &zi64
					z = &zi64
				} else {
					a = zi64
					z = zi64
				}
			}

		} else if strings.Contains(ty.String(), "string") {
			n, ok := convertString(d.Value)
			if ok {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &n
					z = &zs
				} else {
					a = n
					z = zs
				}
			} else {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &zs
					z = &zs
				} else {
					a = zs
					z = zs
				}
			}

		} else if strings.Contains(ty.String(), "int") {
			n, ok := convertInt(d.Value)
			if ok {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &n
					z = &zi
				} else {
					a = n
					z = zi
				}
			} else {
				if v.Field(i).Kind() == reflect.Ptr {
					a = &zi
					z = &zi
				} else {
					a = zi
					z = zi
				}
			}

		} else {
			continue
		}

		if f.Tag.Get("json") == d.Param {
			v.Field(i).Set(reflect.ValueOf(a))
		} else {
			v.Field(i).Set(reflect.ValueOf(z))
		}

	}

	return nil
}

func convertString(v interface{}) (string, bool) {
	s, ok := v.(string)
	if ok {
		return s, true
	}

	n, ok := v.(int)
	if ok {
		return fmt.Sprint(n), true
	}

	return "", false
}

func convertInt(v interface{}) (int, bool) {
	f, ok := v.(float64)
	if ok {
		return int(f), true
	}

	f32, ok := v.(float32)
	if ok {
		return int(f32), true
	}

	i, ok := v.(int)
	if ok {
		return i, true
	}

	s, ok := v.(string)
	if ok {
		i, e := strconv.Atoi(s)
		if e == nil {
			return i, true
		}
	}

	return 0, false
}

func convertInt64(v interface{}) (int64, bool) {
	f, ok := v.(float64)
	if ok {
		return int64(f), true
	}

	f32, ok := v.(float32)
	if ok {
		return int64(f32), true
	}

	i64, ok := v.(int64)
	if ok {
		return i64, true
	}

	i, ok := v.(int)
	if ok {
		return int64(i), true
	}

	s, ok := v.(string)
	if ok {
		i, e := strconv.Atoi(s)
		if e == nil {
			return int64(i), true
		}
	}

	return 0, false
}
