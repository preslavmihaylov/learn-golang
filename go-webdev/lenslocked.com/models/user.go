package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/rand"

	// preload postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var userPwPepper = "some-pepper"
var hmacSecretKey = "secret-hmac-key"

type User struct {
	gorm.Model
	Name              string
	Email             string `gorm:"not null;unique_index"`
	Password          string `gorm:"-"`
	PasswordHash      string `gorm:"not null"`
	RememberToken     string `gorm:"-"`
	RememberTokenHash string `gorm:"not null;unique_index"`
}

func (u *User) GenerateToken() error {
	if u.RememberToken != "" {
		return nil
	}

	tok, err := rand.RememberToken()
	if err != nil {
		return fmt.Errorf("failed to generate new remember token: %s", err)
	}

	u.RememberToken = tok

	return nil
}

func (u *User) isPasswordCorrect(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(userPwPepper+password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, fmt.Errorf("failed comparing password hashes: %s", err)
		}
	}

	return true, nil
}
