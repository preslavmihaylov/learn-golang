package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

type OAuth struct {
	gorm.Model
	oauth2.Token
	UserID  uint   `gorm:"not null;unique_index:user_id_service"`
	Service string `gorm:"not null;unique_index:user_id_service"`
}
