package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/hash"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/rand"

	// preload postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound      = errors.New("models: resource not found")
	ErrInvalidID     = errors.New("models: invalid id")
	ErrUserNotFound  = errors.New("models: user not found")
	ErrWrongPassword = errors.New("models: wrong password")
)

var userPwPepper = "some-pepper"
var hmacSecretKey = "secret-hmac-key"

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

func (u *User) GenerateToken() error {
	if u.Remember != "" {
		return nil
	}

	tok, err := rand.RememberToken()
	if err != nil {
		return fmt.Errorf("failed to generate new remember token: %s", err)
	}

	u.Remember = tok

	return nil
}

func (u *User) hashPassword() error {
	phashBytes, err := bcrypt.GenerateFromPassword([]byte(userPwPepper+u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash user password: %s", err)
	}

	u.Password = ""
	u.PasswordHash = string(phashBytes)

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

func (u *User) hashRememberToken(hmac hash.HMAC) error {
	if u.Remember == "" {
		return nil
	}

	var err error
	u.RememberHash, err = hmac.Hash(u.Remember)
	if err != nil {
		return err
	}

	return nil
}
