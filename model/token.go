package model

import (
	"time"
)

type Token struct {
	Value    string `gorm:"primary_key;unique_index;type:varchar(128)" validate:"required"`
	LastUsed *time.Time
}
