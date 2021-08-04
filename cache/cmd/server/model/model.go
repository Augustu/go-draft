package model

import (
	"time"

	"gorm.io/gorm"
)

type Stats struct {
	gorm.Model
	A          int       `json:"a"`
	B          int       `json:"b"`
	C          int       `json:"c"`
	D          int       `json:"d"`
	OccurredAt time.Time `json:"occurred_at"`
}
