package model

import (
	"time"
)

type Token struct {
	ID       uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"-"`
	// The token
	Value    string `gorm:"primary_key;unique_index;type:varchar(128)"`
	// Created by who
	UserID   uint   `gorm:"primary_key;unique_index;type:varchar(128)" json:"-"`
	CreatedBy User
	// Its last use
	LastUse *time.Time
}
