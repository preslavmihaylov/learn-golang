package context

import (
	"context"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	temp := ctx.Value(userKey)
	if temp == nil {
		return nil
	}

	usr, ok := temp.(*models.User)
	if !ok {
		return nil
	}

	return usr
}
