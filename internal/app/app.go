package app

import (
	"console-chat/internal/app/valid"
	"console-chat/internal/model"
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type app struct {
	UserRepo
}

func (a *app) RegisterUser(ctx context.Context, nickname, password string) (model.User, error) {

	// check if nickname and password are both valid
	if !valid.IsValidNickname(nickname) {
		return model.User{}, model.UserInvalidNickname
	}
	if !valid.IsValidPassword(password) {
		return model.User{}, model.UserInvalidPassword
	}

	var usr model.User
	usr.Nickname = nickname

	// creating hash sum of the password
	hash := sha256.New()
	hash.Write([]byte(password))
	hashSum := hash.Sum(nil)
	usr.HashedPassword = hex.EncodeToString(hashSum)

	return a.AddUser(ctx, usr)
}

func (a *app) SignInUser(ctx context.Context, nickname, password string) (model.User, error) {
	var usr model.User
	var err error

	// getting user with given nickname from repo
	if usr, err = a.GetUser(ctx, nickname); err != nil {
		return model.User{}, err
	}

	// checking password
	hash := sha256.New()
	hash.Write([]byte(password))
	hashSum := hash.Sum(nil)
	if usr.HashedPassword != hex.EncodeToString(hashSum) {
		return model.User{}, model.UserWrongPassword
	} else {
		return usr, nil
	}
}
