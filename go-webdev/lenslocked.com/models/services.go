package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Services struct {
	Gallery GalleryService
	Image   ImageService
	User    UserService
	OAuth   OAuthService
	db      *gorm.DB
}

type ServicesConfigFunc func(*Services) error

func NewServices(cfgFuncs ...ServicesConfigFunc) (*Services, error) {
	var s Services
	for _, cfgFunc := range cfgFuncs {
		if err := cfgFunc(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func WithGorm(dialect, connectionInfo string) ServicesConfigFunc {
	return func(s *Services) error {
		var err error
		s.db, err = gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}

		return nil
	}
}

func WithLogMode(mode bool) ServicesConfigFunc {
	return func(s *Services) error {
		if s.db == nil {
			return fmt.Errorf("database not set")
		}

		s.db.LogMode(mode)

		return nil
	}
}

func WithUserService(pepper, hmacKey string) ServicesConfigFunc {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)

		return nil
	}
}

func WithGalleryService() ServicesConfigFunc {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)

		return nil
	}
}

func WithImageService() ServicesConfigFunc {
	return func(s *Services) error {
		s.Image = NewImageService()

		return nil
	}
}

func WithOAuthService() ServicesConfigFunc {
	return func(s *Services) error {
		s.OAuth = NewOAuthService(s.db)

		return nil
	}
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &pwReset{}, &OAuth{}).Error
}

func (s *Services) DestructiveReset() error {
	s.db.DropTableIfExists(&User{}, &Gallery{}, &pwReset{}, &OAuth{})

	err := s.db.Error
	if err != nil {
		return fmt.Errorf("received error while recreating users table: %s", err)
	}

	return s.AutoMigrate()
}
