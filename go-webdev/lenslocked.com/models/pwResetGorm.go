package models

import "github.com/jinzhu/gorm"

type pwResetDB interface {
	ByToken(token string) (*pwReset, error)
	Create(pwr *pwReset) error
	Delete(id uint) error
}

type pwResetGorm struct {
	db *gorm.DB
}

func (pwrg *pwResetGorm) ByToken(token string) (*pwReset, error) {
	var pwr pwReset
	err := first(pwrg.db.Where("token_hash = ?", token), &pwr)
	if err != nil {
		return nil, err
	}

	return &pwr, nil
}

func (pwrg *pwResetGorm) Create(pwr *pwReset) error {
	return pwrg.db.Create(pwr).Error
}

func (pwrg *pwResetGorm) Delete(id uint) error {
	pwr := pwReset{Model: gorm.Model{ID: id}}

	return pwrg.db.Delete(&pwr).Error
}
