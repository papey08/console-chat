package app

import (
	"console-chat/internal/model"
	"context"
)

type App interface {
	RegisterUser(ctx context.Context, nickname, password string) (model.User, error)
	SignInUser(ctx context.Context, nickname, password string) (model.User, error)
}

type UserRepo interface {
	AddUser(ctx context.Context, u model.User) (model.User, error)
	GetUser(ctx context.Context, nickname string) (model.User, error)
}

func New(repo UserRepo) App {
	return &MyApp{
		UserRepo: repo,
	}
}
