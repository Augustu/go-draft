package gorm

import (
	"gorm.io/driver/mysql"
)

func main() {
	// dsn := ""
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	mysql.New(mysql.Config{})
}
