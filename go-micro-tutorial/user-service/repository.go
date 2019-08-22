package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
)

type userRepository struct {
	db *gorm.DB
}

func (repo *userRepository) GetAll() ([]*proto.User, error) {
	var usrs []*proto.User
	err := repo.db.Find(&usrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users from db: %s", err)
	}

	return usrs, nil
}

func (repo *userRepository) Get(id string) (*proto.User, error) {
	var usr *proto.User
	usr.Id = id
	err := repo.db.First(&usr).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user from db: %s", err)
	}

	return usr, nil
}

func (repo *userRepository) Create(usr *proto.User) error {
	log.Printf("Creating user: %s", usr)
	err := repo.db.Create(&usr).Error
	if err != nil {
		return fmt.Errorf("failed to create user in db: %s", err)
	}

	return nil
}

func (repo *userRepository) GetByEmail(email string) (*proto.User, error) {
	usr := &proto.User{}
	usr.Email = email
	err := repo.db.First(&usr).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email and password from db: %s", err)
	}

	return usr, nil
}
