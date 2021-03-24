package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "github.com/Augustu/go-draft/gorm/naming"
	// mt "github.com/Augustu/go-draft/time"
)

type Event struct {
	gorm.Model
	App     string
	Key     string
	EventID string  `gorm:"index"`
	UID     string  `gorm:"index:idx_member"`
	Score   float64 `gorm:"index:idx_member"`
	OccurAt time.Time
	Remark  string
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/rank?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// using table suffix naming
		// NamingStrategy: naming.TimeNamingStrategy{
		// 	TableSuffix: "_" + mt.String(),
		// },
	})
	if err != nil {
		log.Fatalf("open database failed: %s", err.Error())
	}

	err = db.AutoMigrate(&Event{})
	if err != nil {
		log.Fatalf("auto migrate failed: %s", err.Error())
	}

	e := Event{
		App:     "app",
		Key:     "key",
		EventID: "event_id",
		UID:     "UID",
		Score:   0.0,
		OccurAt: time.Now(),
		Remark:  "remark",
	}

	var res *gorm.DB

	res = db.Create(&e)
	if res.Error != nil {
		log.Printf("insert failed: %s\n", res.Error.Error())
	}

	e.ID = 2
	e.EventID = "event_id_2"
	res = db.Create(&e)
	if res.Error != nil {
		log.Printf("insert failed: %s\n", res.Error.Error())
	}
	log.Printf("insert done: %d", e.ID)

	var events []Event
	r := db.Find(&events)
	if r.Error != nil {
		log.Fatalf("get events failed: %s", r.Error.Error())
	}

	for idx, e := range events {
		fmt.Printf("%d %d %v\n", idx, e.ID, e)
	}

}
