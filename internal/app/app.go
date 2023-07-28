package app

import (
	"console-chat/internal/model"
	"net"
	"net/http"
)

type UserRepo interface {
	CreateUser(u model.User) (error, model.User)
	GetUserByID(id uint64) (error, model.User)
	UpdateUserFields(id uint64, u model.User) (error, model.User)
	DeleteUser(id uint64) error
}

type App interface {
	AddConnection(conn net.Conn) uint64
	Chat(w http.ResponseWriter, r *http.Request)
}
