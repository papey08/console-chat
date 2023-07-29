package model

import "errors"

var UserNotFound = errors.New("could not find required user")
var UserRepoError = errors.New("something wrong with user repo")
var UserAlreadyExists = errors.New("user with required nickname already exists")
var UserWrongPassword = errors.New("wrong password of required user")
var UserInvalidNickname = errors.New("user has invalid nickname")
var UserInvalidPassword = errors.New("user has invalid password")
