package models

import (
	"strings"
)

var (
	ErrNotFound modelError = "models: resource not found"
	ErrNilModel modelError = "models: given resource is nil"

	// User Errors
	ErrIDInvalid             modelError = "models: invalid id"
	ErrUserNotFound          modelError = "models: user not found"
	ErrPasswordWrong         modelError = "models: wrong password"
	ErrPasswordTooShort      modelError = "models: password too short. Must be at least 8 characters long"
	ErrPasswordRequired      modelError = "models: password is required"
	ErrEmailRequired         modelError = "models: email required"
	ErrEmailInvalid          modelError = "models: invalid email"
	ErrEmailTaken            modelError = "models: email is taken"
	ErrRememberTokenTooShort modelError = "models: remember token is too short"
	ErrRememberTokenRequired modelError = "models: remember token is required"
	ErrRememberTokenInvalid  modelError = "models: remember token is invalid"

	// Gallery Errors
	ErrUserIDRequired modelError = "models: user ID is required"
	ErrTitleRequired  modelError = "models: title is required"

	// OAuth Errors
	ErrOAuthServiceRequired modelError = "models: oauth service is required"
	ErrOAuthTokenRequired   modelError = "models: oauth token is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])

	return strings.Join(split, " ")
}
