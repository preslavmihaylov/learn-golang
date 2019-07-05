package models

import "github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/rand"

type pwResetValidationFunc func(pwr *pwReset) error

func runPwResetValidationFuncs(pwr *pwReset, fns ...pwResetValidationFunc) error {
	for _, fn := range fns {
		err := fn(pwr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pwrv *pwResetValidator) requireUserID(pwr *pwReset) error {
	if pwr.UserID <= 0 {
		return ErrUserIDRequired
	}

	return nil
}

func (pwrv *pwResetValidator) setTokenIfUnset(pwr *pwReset) error {
	if pwr.Token != "" {
		return nil
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}

	pwr.Token = token
	return nil
}

func (pwrv *pwResetValidator) hmacToken(pwr *pwReset) error {
	if pwr.Token == "" {
		return nil
	}

	var err error
	pwr.TokenHash, err = pwrv.hmac.Hash(pwr.Token)
	if err != nil {
		return err
	}

	return nil
}
