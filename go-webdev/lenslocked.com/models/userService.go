package models

import (
	"github.com/jinzhu/gorm"
)

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
}

func NewUserService(db *gorm.DB) UserService {
	ug := userGorm{db}
	return &userService{
		UserDB: NewUserValidator(&ug),
	}
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	usr, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	isPassCorrect, err := usr.isPasswordCorrect(password)
	if err != nil {
		return nil, err
	}

	if !isPassCorrect {
		return nil, ErrPasswordWrong
	}

	return usr, nil
}
