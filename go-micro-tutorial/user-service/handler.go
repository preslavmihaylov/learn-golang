package main

import (
	"context"
	"fmt"
	"log"

	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo         *userRepository
	tokenService Authable
}

func (us *userService) Create(ctx context.Context, usr *proto.User, resp *proto.Response) error {
	hPass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err)
	}

	usr.Password = string(hPass)
	err = us.repo.Create(usr)
	if err != nil {
		return fmt.Errorf("failed to create user: %s", err)
	}

	resp.User = usr
	return nil
}

func (us *userService) Get(ctx context.Context, req *proto.User, resp *proto.Response) error {
	usr, err := us.repo.Get(req.Id)
	if err != nil {
		return fmt.Errorf("failed to get user: %s", err)
	}

	resp.User = usr
	return nil
}

func (us *userService) GetAll(ctx context.Context, req *proto.Request, resp *proto.Response) error {
	usrs, err := us.repo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all users: %s", err)
	}

	resp.Users = usrs
	return nil
}

func (us *userService) Auth(ctx context.Context, req *proto.User, tok *proto.Token) error {
	log.Printf("Authenticating <%s> (%s)...", req.Email, req.Password)
	usr, err := us.repo.GetByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("user does not exist: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password))
	if err != nil {
		return fmt.Errorf("password mismatch: %s", err)
	}

	tok.Token, err = us.tokenService.Encode(usr)
	if err != nil {
		return fmt.Errorf("failed to encode token: %s", err)
	}

	return nil
}

func (us *userService) ValidateToken(ctx context.Context, inTok *proto.Token, outTok *proto.Token) error {
	return nil
}
