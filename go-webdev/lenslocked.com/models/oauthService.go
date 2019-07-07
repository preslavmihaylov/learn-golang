package models

import "github.com/jinzhu/gorm"

type OAuthService interface {
	OAuthDB
}

func NewOAuthService(db *gorm.DB) OAuthService {
	return &oauthValidator{
		OAuthDB: &oauthGorm{db: db},
	}
}
