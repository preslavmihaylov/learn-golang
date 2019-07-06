package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/hash"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Authenticate(email, password string) (*User, error)
	InitiatePasswordReset(email string) (string, error)
	CompletePasswordReset(token, newPassword string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
	pepper    string
	pwResetDB pwResetDB
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := userGorm{db}
	hmac := hash.NewHMAC(hmacKey)
	return &userService{
		UserDB:    newUserValidator(&ug, pepper, hmac),
		pepper:    pepper,
		pwResetDB: newPwResetValidator(&pwResetGorm{db}, hmac),
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

func (us *userService) InitiatePasswordReset(email string) (string, error) {
	usr, err := us.ByEmail(email)
	if err != nil {
		return "", err
	}

	pwr := newPwReset(usr.ID)
	err = us.pwResetDB.Create(pwr)
	if err != nil {
		return "", err
	}

	return pwr.Token, nil
}

func (us *userService) CompletePasswordReset(token, newPassword string) (*User, error) {
	pwr, err := us.pwResetDB.ByToken(token)
	if err != nil {
		switch err {
		case ErrNotFound:
			return nil, ErrRememberTokenInvalid
		default:
			return nil, err
		}
	}

	if time.Now().Sub(pwr.CreatedAt) > time.Hour*12 {
		return nil, ErrRememberTokenInvalid
	}

	usr, err := us.ByID(pwr.UserID)
	if err != nil {
		return nil, err
	}

	usr.Password = newPassword
	err = us.Update(usr)
	if err != nil {
		return nil, err
	}

	err = us.pwResetDB.Delete(pwr.ID)
	if err != nil {
		return nil, err
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
