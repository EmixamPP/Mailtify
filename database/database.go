package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"mailtify/model"
)

const DATABASE_ERROR_MESSAGE = "internal error"

// New creates a GormDB.
// It returns nil and an error if a problem occurs.
func New(dialect, connection string, tokenSize uint8) (*GormDB, error) {
	db, err := gorm.Open(dialect, connection)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(new(model.Token), new(model.Token)).Error; err != nil {
		return nil, err
	}

	return &GormDB{db: db, TokenSize: tokenSize}, nil
}

// GormDB is a wrapper for the gorm framework and other needed parameters.
type GormDB struct {
	db        *gorm.DB
	TokenSize uint8
}

// Close closes the database connection.
func (d *GormDB) Close() {
	d.db.Close()
}
