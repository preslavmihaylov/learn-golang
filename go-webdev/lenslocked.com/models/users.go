package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	// preload postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: invalid id")
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique_index"`
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(connInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %s", err)
	}

	db.LogMode(true)
	return &UserService{db: db}, nil
}

func (us *UserService) AutoMigrate() error {
	err := us.db.AutoMigrate(&User{}).Error
	if err != nil {
		return fmt.Errorf("failed to setup auto-migrate on users table: %s", err)
	}

	return nil
}

func (us *UserService) ByID(id uint) (*User, error) {
	var u User
	err := first(us.db.Where("id = ?", id), &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var u User
	err := first(us.db.Where("email = ?", email), &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *UserService) Create(u *User) error {
	err := us.db.Create(u).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %s", err)
	}

	return nil
}

func (us *UserService) Update(u *User) error {
	err := us.db.Save(u).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	delUsr := User{Model: gorm.Model{ID: id}}
	err := us.db.Delete(&delUsr).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %s", err)
	}

	return nil
}

func (us *UserService) Close() error {
	return us.db.Close()
}

func (us *UserService) DestructiveReset() error {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})

	err := us.db.Error
	if err != nil {
		return fmt.Errorf("received error while recreating users table: %s", err)
	}

	return nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("unexpected error while querying db: %s", err)
		}
	}

	return nil
}
