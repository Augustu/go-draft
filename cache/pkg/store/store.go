package store

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Augustu/go-draft/cache/pkg/config"
)

type Store struct {
	DB *gorm.DB
}

func New(c *config.Store) *Store {
	db, err := gorm.Open(mysql.Open(c.URI), nil)
	if err != nil {
		return nil
	}

	return &Store{
		DB: db,
	}
}
