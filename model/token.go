package model

import (
	"time"
)

type Token struct {
	Value    string     `gorm:"primary_key;unique_index;type:varchar(128)" validate:"required" json:"value"`
	LastUsed *time.Time `json:"lastused"`
}
