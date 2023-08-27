package model

import (
	"time"
)

type User struct {
	ID             uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Username       string `gorm:"unique_index;type:varchar(128)"`
	Password       string
	Tokens         []Token
	LastConnection *time.Time
}
