package ginserver

import (
	"console-chat/internal/model"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	Nickname       string `json:"nickname"`
	HashedPassword string `json:"hashed_password"`
}

func getUserResponse(usr model.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			Nickname:       usr.Nickame,
			HashedPassword: usr.HashedPassword,
		},
		"error": nil,
	}
}

func postUserResponse(usr model.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			Nickname:       usr.Nickame,
			HashedPassword: usr.HashedPassword,
		},
		"error": nil,
	}
}

func ErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
