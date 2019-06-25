package models

import (
	"fmt"
	"strings"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/rand"
	"golang.org/x/crypto/bcrypt"
)

type userValidationFunc func(u *User) error

func runUserValidationFuncs(u *User, fns ...userValidationFunc) error {
	for _, fn := range fns {
		err := fn(u)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uv *userValidator) hashPassword(u *User) error {
	if u.Password == "" {
		return nil
	}

	phashBytes, err := bcrypt.GenerateFromPassword([]byte(userPwPepper+u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash user password: %s", err)
	}

	u.Password = ""
	u.PasswordHash = string(phashBytes)

	return nil
}

func (uv *userValidator) requirePassword(u *User) error {
	if u.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) requirePasswordHash(u *User) error {
	if u.PasswordHash == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) passwordMinLength(minLength int) userValidationFunc {
	return userValidationFunc(func(u *User) error {
		if u.Password == "" {
			return nil
		}

		if len(u.Password) < minLength {
			return ErrPasswordTooShort
		}

		return nil
	})
}

func (uv *userValidator) createRememberToken(u *User) error {
	if u.RememberToken != "" {
		return nil
	}

	tok, err := rand.RememberToken()
	if err != nil {
		return fmt.Errorf("failed to create remember token for user: %s", err)
	}

	u.RememberToken = tok

	return nil
}

func (uv *userValidator) hashRememberToken(u *User) error {
	if u.RememberToken == "" {
		return nil
	}

	var err error
	u.RememberTokenHash, err = uv.hmac.Hash(u.RememberToken)
	if err != nil {
		return err
	}

	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValidationFunc {
	return userValidationFunc(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}

		return nil
	})
}

func (uv *userValidator) normalizeEmail(u *User) error {
	u.Email = strings.TrimSpace(u.Email)
	u.Email = strings.ToLower(u.Email)

	return nil
}

func (uv *userValidator) requireEmail(u *User) error {
	if u.Email == "" {
		return ErrEmailRequired
	}

	return nil
}

func (uv *userValidator) emailIsValid(u *User) error {
	if u.Email == "" {
		return nil
	}

	if !uv.emailRegex.MatchString(u.Email) {
		return ErrEmailInvalid
	}

	return nil
}

func (uv *userValidator) emailIsNotTaken(u *User) error {
	foundUser, err := uv.ByEmail(u.Email)
	if err != nil {
		switch err {
		case ErrNotFound:
			return nil
		default:
			return err
		}
	}

	if foundUser.ID != u.ID {
		return ErrEmailTaken
	}

	return nil
}

func (uv *userValidator) rememberTokenMinBytes(minBytes int) userValidationFunc {
	return userValidationFunc(func(u *User) error {
		if u.RememberToken == "" {
			return nil
		}

		tokenLen, err := rand.NBytes(u.RememberToken)
		if err != nil {
			return fmt.Errorf("failed to get length of remember token: %s", err)
		}

		if tokenLen < minBytes {
			return ErrRememberTokenTooShort
		}

		return nil
	})
}

func (uv *userValidator) requireRememberHash(u *User) error {
	if u.RememberTokenHash == "" {
		return ErrRememberTokenRequired
	}

	return nil
}
