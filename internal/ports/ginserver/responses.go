package ginserver

import (
	"console-chat/internal/model"

	"github.com/gin-gonic/gin"
)

type tokenResponse struct {
	TokenString string `json:"token_string"`
}

func getUserResponse(token string) *gin.H {
	return &gin.H{
		"data": tokenResponse{
			TokenString: token,
		},
		"error": nil,
	}
}

type userResponse struct {
	Nickname       string `json:"nickname"`
	HashedPassword string `json:"hashed_password"`
}

func postUserResponse(usr model.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			Nickname:       usr.Nickname,
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
