package models

import (
	"fmt"
)

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
}

func NewUserService(connInfo string) (UserService, error) {
	ug, err := newUserGorm(connInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate user gorm: %s", err)
	}

	return &userService{
		UserDB: NewUserValidator(ug),
	}, nil
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
