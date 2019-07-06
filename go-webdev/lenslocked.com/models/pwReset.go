package models

import "github.com/jinzhu/gorm"

type pwReset struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"-"`
	TokenHash string `gorm:"not null;unique_index"`
}

func newPwReset(userID uint) *pwReset {
	return &pwReset{
		UserID: userID,
	}
}
