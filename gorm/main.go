package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	me "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "github.com/Augustu/go-draft/gorm/naming"
	// mt "github.com/Augustu/go-draft/time"
)

type Event struct {
	gorm.Model
	App       string
	Key       string
	EventID   string    `gorm:"index"`
	UID       string    `gorm:"index"`
	Score     float64   `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
	OccurAt   time.Time
	Remark    string
}

func test() {
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

func RandomInt(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(end - start)
	random = start + random
	return random
}

func RandomFloat64() float64 {
	return rand.Float64()
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

type client struct {
	dsn string
	db  *gorm.DB
}

func (c *client) init() {
	db, err := gorm.Open(mysql.Open(c.dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&Event{})
	c.db = db
}

func (c *client) Create(n int) {

	batchSize := 1000

	for i := 0; i < n; i = i + batchSize {

		var events []Event

		for j := 0; j < batchSize; j++ {
			events = append(events, Event{
				App:       RandString(8),
				Key:       RandString(12),
				EventID:   RandString(16),
				UID:       RandString(12),
				Score:     RandomFloat64() * 100,
				CreatedAt: time.Now(),
				OccurAt:   time.Now(),
				Remark:    RandString(32),
			})
		}

		tx := c.db.CreateInBatches(events, batchSize)
		// tx := c.db.Create(&events)
		if tx.Error != nil {
			log.Fatal(tx.Error)
		}

		log.Printf("write %d batch", i+batchSize)
	}
}

func (c *client) QueryPage(offset, limit int) {
	var events []Event
	tx := c.db.Offset(offset).Limit(limit).Order("score DESC").Find(&events)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	for idx, e := range events {
		fmt.Println(idx, e)
	}
}

func (c *client) QueryUID(offset, limit int, uid string) {
	var events []Event
	tx := c.db.Offset(offset).Limit(limit).Where("uid = ?", "IRKXCPGIYMPG").Find(&events)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	for idx, e := range events {
		fmt.Println(idx, e)
	}
}

func (c *client) Delete() {

	const TIME_LAYOUT = "2006-01-02 15:04:05"

	t, err := time.Parse(TIME_LAYOUT, "2018-09-10 00:00:00")
	if err != nil {
		log.Fatal(err)
	}

	e := Event{
		App:     "app",
		Key:     "key",
		EventID: "event_id",
		UID:     "UID",
		Score:   100.0,
		OccurAt: t,
		Remark:  "remark",
	}

	tx := c.db.Create(&e)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	tx = c.db.Delete(&e)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	// ne := Event{
	// 	App:     "app",
	// 	Key:     "key",
	// 	EventID: "event_id",
	// 	UID:     "UID",
	// 	Score:   100.0,
	// 	OccurAt: t,
	// 	Remark:  "remark",
	// }
	tx = c.db.Create(&e)
	if tx.Error != nil {
		if tx.Error.(*me.MySQLError).Number == 1062 {
			log.Printf("duplicate key: %v", e)
		}

		log.Fatal(tx.Error)
	}
}

func main1() {
	fmt.Println(RandString(32))
	fmt.Println(RandomFloat64())
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/rank?charset=utf8mb4&parseTime=True&loc=Local"

	c := &client{dsn: dsn}
	c.init()

	// c.db.Create(&Event{
	// 	App:       RandString(8),
	// 	Key:       RandString(12),
	// 	EventID:   RandString(16),
	// 	UID:       RandString(12),
	// 	Score:     RandomFloat64() * 100,
	// 	CreatedAt: time.Now(),
	// 	OccurAt:   time.Now(),
	// 	Remark:    RandString(32),
	// })

	// c.Create(10000)

	// query page
	// c.QueryPage(0, 5)
	// c.QueryPage(5, 5)

	// QueryUID by uid
	// c.QueryUID(0, 5, "IRKXCPGIYMPG")

	c.Delete()

}
