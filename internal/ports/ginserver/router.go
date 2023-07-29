package ginserver

import (
	"console-chat/internal/app"
	"console-chat/internal/ports/wsserver"

	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, ws *wsserver.WsServer, a app.App) {
	r.GET("/chat", gin.WrapF(ws.Chat))
	r.GET("/users/:user_nickname", getUser(a))
	r.POST("users", postUser(a))
}
