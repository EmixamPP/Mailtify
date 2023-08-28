package model

import (
	"time"
)

type User struct {
	ID             uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"-"`
	Username       string `gorm:"unique_index;type:varchar(128)"`
	Password       string `json:"-"`
	Tokens         []Token
	LastConnection *time.Time
	Admin          bool
}
