package model

import (
	"time"
)

type Token struct {
	ID       uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"-"`
	Value    string `gorm:"primary_key;unique_index;type:varchar(128)"`
	UserID   uint   `gorm:"primary_key;unique_index;type:varchar(128)" json:"-"`
	CreatedBy User
	LastUse *time.Time
}
