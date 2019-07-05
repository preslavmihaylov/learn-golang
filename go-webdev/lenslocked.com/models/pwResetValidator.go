package models

import "github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/hash"

type pwResetValidator struct {
	pwResetDB
	hmac hash.HMAC
}

func newPwResetValidator(db pwResetDB, hmac hash.HMAC) *pwResetValidator {
	return &pwResetValidator{
		pwResetDB: db,
		hmac:      hmac,
	}
}

func (pwrv *pwResetValidator) ByToken(token string) (*pwReset, error) {
	pwr := pwReset{Token: token}
	err := runPwResetValidationFuncs(&pwr, pwrv.hmacToken)
	if err != nil {
		return nil, err
	}

	return pwrv.pwResetDB.ByToken(pwr.TokenHash)
}

func (pwrv *pwResetValidator) Create(pwr *pwReset) error {
	err := runPwResetValidationFuncs(pwr, pwrv.requireUserID, pwrv.setTokenIfUnset, pwrv.hmacToken)
	if err != nil {
		return err
	}

	return pwrv.pwResetDB.Create(pwr)
}

func (pwrv *pwResetValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}

	return pwrv.pwResetDB.Delete(id)
}
