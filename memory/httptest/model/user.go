package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string
	Password    string
	Gender      string
	Age         int
	Title       string
	Description string
}
