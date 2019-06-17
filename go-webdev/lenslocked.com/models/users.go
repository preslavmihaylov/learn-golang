package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	// preload postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound      = errors.New("models: resource not found")
	ErrInvalidID     = errors.New("models: invalid id")
	ErrUserNotFound  = errors.New("models: user not found")
	ErrWrongPassword = errors.New("models: wrong password")
)

var userPwPepper = "some-pepper"

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
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
	err := u.hashPassword()
	if err != nil {
		return err
	}

	err = us.db.Create(u).Error
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

func (u *User) hashPassword() error {
	phashBytes, err := bcrypt.GenerateFromPassword([]byte(userPwPepper+u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash user password: %s", err)
	}

	u.Password = ""
	u.PasswordHash = string(phashBytes)

	return nil
}

func (u *User) isPasswordCorrect(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(userPwPepper+password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, fmt.Errorf("failed comparing password hashes: %s", err)
		}
	}

	return true, nil
}
