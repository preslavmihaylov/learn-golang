package models

import "errors"

var (
	ErrIDInvalid             = errors.New("models: invalid id")
	ErrUserNotFound          = errors.New("models: user not found")
	ErrPasswordWrong         = errors.New("models: wrong password")
	ErrPasswordTooShort      = errors.New("models: password too short. Must be at least 8 characters long")
	ErrPasswordRequired      = errors.New("models: password is required")
	ErrEmailRequired         = errors.New("models: email required")
	ErrEmailInvalid          = errors.New("models: invalid email")
	ErrEmailTaken            = errors.New("models: email is taken")
	ErrRememberTokenTooShort = errors.New("models: remember token is too short")
	ErrRememberTokenRequired = errors.New("models: remember token is required")
)
