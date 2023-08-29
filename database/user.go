package database

import (
	"github.com/jinzhu/gorm"

	"mailtify/model"
)

// CreateUser creates an user in the database
// An error are returned if a problem has occurred.
func (d *GormDB) CreateUser(username, password string, admin bool) error {
	user := model.User{Username: username, Password: password, Admin: admin} // TODO store hash
	return d.db.Create(&user).Error
}

// Deluser deletes an user from the database.
// An error is returned if a problem has occurred.
func (d *GormDB) DelUser(user *model.User) error {
	return d.db.Delete(&user).Error
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
	if err := d.db.Model(&user).Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (d *GormDB) IsAdminExists() (bool, error) {
	var count int
	result := d.db.Model(&model.User{}).Where("admin = ?", true).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
