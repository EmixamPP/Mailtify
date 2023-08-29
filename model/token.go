package model

import (
	"time"
)

type Token struct {
	ID      uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"-"`
	Value   string `gorm:"unique_index;type:varchar(128)"`
	UserID  uint   `json:"-"`
	User    User
	LastUse *time.Time
}
