package models

import (
	"regexp"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/hash"
)

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	pepper     string
	emailRegex *regexp.Regexp
}

func NewUserValidator(userDB UserDB, pepper, hmacKey string) *userValidator {
	return &userValidator{
		UserDB:     userDB,
		hmac:       hash.NewHMAC(hmacKey),
		pepper:     pepper,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	var u User
	u.Email = email

	err := runUserValidationFuncs(&u, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(u.Email)
}

func (uv *userValidator) ByRememberToken(token string) (*User, error) {
	u := User{RememberToken: token}
	err := runUserValidationFuncs(&u, uv.hashRememberToken)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByRememberToken(u.RememberTokenHash)
}

func (uv *userValidator) Create(u *User) error {
	err := runUserValidationFuncs(u,
		uv.normalizeEmail, uv.requireEmail, uv.emailIsValid, uv.emailIsNotTaken,
		uv.passwordMinLength(8), uv.requirePassword, uv.hashPassword, uv.requirePasswordHash,
		uv.createRememberToken, uv.rememberTokenMinBytes(32), uv.hashRememberToken, uv.requireRememberHash)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(u)
}

func (uv *userValidator) Update(u *User) error {
	err := runUserValidationFuncs(u,
		uv.normalizeEmail, uv.requireEmail, uv.emailIsValid, uv.emailIsNotTaken,
		uv.passwordMinLength(8), uv.hashPassword, uv.requirePasswordHash,
		uv.rememberTokenMinBytes(32), uv.hashRememberToken, uv.requireRememberHash)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(u)
}

func (uv *userValidator) Delete(id uint) error {
	var u User
	u.ID = id

	err := runUserValidationFuncs(&u, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}
