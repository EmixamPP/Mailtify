package model

type User struct {
	ID       uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT" json:"-"`
	Username string `gorm:"unique_index;type:varchar(128)"`
	Password string `json:"-"`
	Admin    bool
}
