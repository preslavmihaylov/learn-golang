package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

func NewServices(connInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate user gorm: %s", err)
	}

	db.LogMode(true)

	return &Services{
		Gallery: NewGalleryService(db),
		User:    NewUserService(db),
		db:      db,
	}, nil
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

func (s *Services) DestructiveReset() error {
	s.db.DropTableIfExists(&User{}, &Gallery{})

	err := s.db.Error
	if err != nil {
		return fmt.Errorf("received error while recreating users table: %s", err)
	}

	return s.AutoMigrate()
}
