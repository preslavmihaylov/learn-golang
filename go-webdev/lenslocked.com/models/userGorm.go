package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRememberToken(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var u User
	err := first(ug.db.Where("id = ?", id), &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var u User
	err := first(ug.db.Where("email = ?", email), &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (ug *userGorm) ByRememberToken(rememberHash string) (*User, error) {
	var u User
	err := first(ug.db.Where("remember_token_hash = ?", rememberHash), &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (ug *userGorm) Create(u *User) error {
	err := ug.db.Create(u).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %s", err)
	}

	return nil
}

func (ug *userGorm) Update(u *User) error {
	err := ug.db.Save(u).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

func (ug *userGorm) Delete(id uint) error {
	delUsr := User{Model: gorm.Model{ID: id}}
	err := ug.db.Delete(&delUsr).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %s", err)
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
