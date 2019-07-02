package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/rand"

	// preload postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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
