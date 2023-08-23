package database

import (
	"crypto/rand"

	"github.com/jinzhu/gorm"

	"mailtify/model"
)

const TOKEN_CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateToken returns a new random token of tokenSize alphenumeric character.
// nil and an error are returned if a problem has occurred.
func generateToken(tokenSize int) (*model.Token, error) {
	token := make([]byte, tokenSize)

	// generate a token of random byte
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}

	// translate into alphenumeric
	for i := 0; i < tokenSize; i++ {
		token[i] = TOKEN_CHARSET[token[i]%byte(len(TOKEN_CHARSET))]
	}

	return &model.Token{Value: string(token)}, nil
}

// NewToken creates an unique token in the storage and returns it.
// nil and an error are returned if a problem has occurred.
func (d *GormDB) NewToken() (*model.Token, error) {
	token, err := generateToken(d.TokenSize)
	if err != nil {
		return nil, err
	}

	err = d.DB.Create(token).Error
	if err != nil {
		return nil, err
	}

	return token, nil
}

// DelToken deletes a token in the storage.
// An error is returned if a problem has occurred.
func (d *GormDB) DelToken(token *model.Token) error {
	return d.DB.Delete(token).Error
}

// IsTokenExists checks if a token is the storage, returns true if so, otherwise false.
// False and an error is returned if a problem has occurred.
func (d *GormDB) IsTokenExists(token *model.Token) (bool, error) {
	if err := d.DB.Where("value = ?", token.Value).First(&model.Token{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
