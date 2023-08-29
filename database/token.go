package database

import (
	"crypto/rand"

	"github.com/jinzhu/gorm"

	"mailtify/model"
)

const TOKEN_CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateToken returns a new random token of tokenSize alphenumeric character.
// nil and an error are returned if a problem has occurred.
func generateToken(tokenSize uint8) (*model.Token, error) {
	token := make([]byte, tokenSize)

	// generate a token of random byte
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}

	// translate into alphenumeric
	for i := uint8(0); i < tokenSize; i++ {
		token[i] = TOKEN_CHARSET[token[i]%byte(len(TOKEN_CHARSET))]
	}

	return &model.Token{Value: string(token)}, nil
}

// NewToken creates an unique token in the database and returns it.
// nil and an error are returned if a problem has occurred.
func (d *GormDB) NewToken(createdBy *model.User) (*model.Token, error) {
	token, err := generateToken(d.TokenSize)
	if err != nil {
		return nil, err
	}

	token.UserID = createdBy.ID
	token.User = *createdBy

	err = d.db.Create(&token).Error
	if err != nil {
		return nil, err
	}

	createdBy.Tokens = append(createdBy.Tokens, *token)

	return token, nil
}

// SaveToken update a token of the database.
// An error are returned if a problem has occurred.
func (d *GormDB) UpdateToken(token *model.Token) error {
	return d.db.Save(&token).Error
}

// DelToken deletes a token from the database.
// An error is returned if a problem has occurred.
func (d *GormDB) DelToken(token *model.Token) error {
	return d.db.Delete(&token).Error
}

// GetToken checks if a token is in the database,
// returns the token model if so, otherwise nil.
// nil and an error is returned if a problem has occurred.
func (d *GormDB) GetToken(value string) (*model.Token, error) {
	var token model.Token
	if err := d.db.Model(&token).Where("value = ?", value).First(&token).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// GetTokens returns all tokens in the dartabase.
// An empty list and an error are returned if a problem has occured.
func (d *GormDB) GetTokens() ([]model.Token, error) {
	var tokens []model.Token
	err := d.db.Preload("User").Find(&tokens).Error
	return tokens, err
}
