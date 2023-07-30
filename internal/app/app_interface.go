package app

import (
	"console-chat/internal/model"
	"context"
)

type App interface {
	// Register user checks nickname and password validity and adds new user to the repo
	RegisterUser(ctx context.Context, nickname, password string) (model.User, error)

	// SignInUser finds user in user repo by nickname and checks if password is right
	SignInUser(ctx context.Context, nickname, password string) (model.User, error)
}

type UserRepo interface {
	// AddUser adds new user to the repo
	AddUser(ctx context.Context, u model.User) (model.User, error)

	// GetUser finds user in the repo by nickname
	GetUser(ctx context.Context, nickname string) (model.User, error)
}

func New(repo UserRepo) App {
	return &app{
		UserRepo: repo,
	}
}
