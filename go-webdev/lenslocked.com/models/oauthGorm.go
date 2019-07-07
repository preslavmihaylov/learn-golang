package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type OAuthDB interface {
	Find(userID uint, service string) (*OAuth, error)
	Create(oauth *OAuth) error
	Delete(id uint) error
}

type oauthGorm struct {
	db *gorm.DB
}

func (og *oauthGorm) Find(userID uint, service string) (*OAuth, error) {
	var oauth OAuth
	err := first(og.db.Where("user_id = ? AND service = ?", userID, service), &oauth)
	if err != nil {
		return nil, err
	}

	return &oauth, nil
}

func (og *oauthGorm) Create(oauth *OAuth) error {
	err := og.db.Create(oauth).Error
	if err != nil {
		return fmt.Errorf("failed to create oauth model: %s", err)
	}

	return nil
}

func (og *oauthGorm) Delete(id uint) error {
	var oauth OAuth
	oauth.ID = id

	err := og.db.Unscoped().Delete(&oauth).Error
	if err != nil {
		return err
	}

	return nil
}
