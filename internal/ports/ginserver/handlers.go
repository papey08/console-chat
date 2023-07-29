package ginserver

import (
	"console-chat/internal/app"
	"console-chat/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		nickname := c.Param("user_nickname")
		var reqBody getUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		usr, getErr := a.SignInUser(c, nickname, reqBody.Password)
		switch getErr {
		case model.UserNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(getErr))
		case model.UserWrongPassword:
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(getErr))
		case nil:
			c.JSON(http.StatusOK, getUserResponse(usr))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(getErr))
		}
	}
}

func postUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody postUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		usr, postErr := a.RegisterUser(c, reqBody.Nickname, reqBody.Password)
		switch postErr {
		case model.UserAlreadyExists:
			c.AbortWithStatusJSON(http.StatusConflict, ErrorResponse(postErr))
		case model.UserInvalidNickname:
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(postErr))
		case model.UserInvalidPassword:
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(postErr))
		case nil:
			c.JSON(http.StatusOK, postUserResponse(usr))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(postErr))
		}
	}
}
