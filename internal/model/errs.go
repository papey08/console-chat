package model

import "errors"

var UserNotFound = errors.New("could not find required user")
var UserRepoError = errors.New("something wrong with user repo")
