package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
	pepper string
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := userGorm{db}
	return &userService{
		UserDB: NewUserValidator(&ug, pepper, hmacKey),
		pepper: pepper,
	}
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	usr, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	isPassCorrect, err := us.isPasswordCorrect(usr, password)
	if err != nil {
		return nil, err
	}

	if !isPassCorrect {
		return nil, ErrPasswordWrong
	}

	return usr, nil
}

func (us *userService) isPasswordCorrect(usr *User, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(us.pepper+password))
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
