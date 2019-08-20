package main

import (
	"context"
	"fmt"

	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
)

type userService struct {
	repo *userRepository
}

func (us *userService) Create(ctx context.Context, usr *proto.User, resp *proto.Response) error {
	err := us.repo.Create(usr)
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
	_, err := us.repo.GetByEmailAndPassword(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate user: %s", err)
	}

	tok.Token = "testingabc"
	return nil
}

func (us *userService) ValidateToken(ctx context.Context, inTok *proto.Token, outTok *proto.Token) error {
	return nil
}
