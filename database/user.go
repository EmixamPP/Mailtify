package database

import (
	"github.com/jinzhu/gorm"

	"mailtify/model"
)

// CreateUser creates an user in the database
// An error are returned if a problem has occurred.
func (d *GormDB) CreateUser(username, password string) error {
	user := model.User{Username: username, Password: password} // TODO store hash
	return d.db.Create(user).Error
}

// Deluser deletes an user from the database.
// An error is returned if a problem has occurred.
func (d *GormDB) DelUser(username string) error {
	return d.db.Where("username = ?", username).Delete(&model.Token{}).Error
}

// GetUsers returns all users in the dartabase.
// An empty list and an error are returned if a problem has occured.
func (d *GormDB) GetUsers() ([]model.User, error) {
	var users []model.User
	err := d.db.Find(&users).Error
	return users, err
}

// GetUser checks if an user is in the database,
// returns the user model if so, otherwise nil.
// nil and an error is returned if a problem has occurred.
func (d *GormDB) GetUser(username, password string) (*model.User, error) {
	var user model.User
	if err := d.db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
