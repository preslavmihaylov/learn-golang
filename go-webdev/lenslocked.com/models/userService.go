package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/hash"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/rand"
)

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

func NewUserService(connInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %s", err)
	}

	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{db: db, hmac: hmac}, nil
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

func (us *UserService) ByRememberToken(token string) (*User, error) {
	rememberHash, err := us.hmac.Hash(token)
	if err != nil {
		return nil, fmt.Errorf("failed to hash remember token: %s", err)
	}

	var u User
	err = first(us.db.Where("remember_hash = ?", rememberHash), &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *UserService) Create(u *User) error {
	err := u.hashPassword()
	if err != nil {
		return err
	}

	if u.Remember == "" {
		tok, err := rand.RememberToken()
		if err != nil {
			return fmt.Errorf("failed to create remember token for user: %s", err)
		}

		u.Remember = tok
	}

	err = u.hashRememberToken(us.hmac)
	if err != nil {
		return fmt.Errorf("failed to hash remember token: %s", err)
	}

	err = us.db.Create(u).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %s", err)
	}

	return nil
}

func (us *UserService) Update(u *User) error {
	err := u.hashRememberToken(us.hmac)
	if err != nil {
		return fmt.Errorf("failed to hash remember token: %s", err)
	}

	err = us.db.Save(u).Error
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

func (us *UserService) Authenticate(email, password string) (*User, error) {
	usr, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	isPassCorrect, err := usr.isPasswordCorrect(password)
	if err != nil {
		return nil, err
	}

	if !isPassCorrect {
		return nil, ErrWrongPassword
	}

	return usr, nil
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
