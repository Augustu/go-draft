package main

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Augustu/go-draft/gorm/naming"
	mt "github.com/Augustu/go-draft/time"
)

type Event struct {
	gorm.Model
	App     string
	Key     string
	EventID string  `gorm:"primaryKey"`
	UID     string  `gorm:"index:idx_member"`
	Score   float64 `gorm:"index:idx_member"`
	OccurAt time.Time
	Remark  string
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/rank?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: naming.TimeNamingStrategy{
			TableSuffix: "_" + mt.String(),
		},
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

	res := db.Create(&e)
	if res.Error != nil {
		log.Println("insert failed")
	}

	log.Printf("insert done: %d", e.ID)
}
